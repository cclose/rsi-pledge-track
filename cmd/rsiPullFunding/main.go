package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.org/cclose/rsi-pledge-track/database"
	"github.org/cclose/rsi-pledge-track/service"

	"github.org/cclose/rsi-pledge-track/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const developerEmail = "pulsar2612@gmail.com"
const rsiPledgeDataURL = "https://robertsspaceindustries.com/api/stats/getCrowdfundStats"

func makeRequest(url string, data *model.RSICrowdFundStatsRequest) (*http.Response, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getPledgeData(timestamp string) (*model.PledgeData, error) {
	data := &model.RSICrowdFundStatsRequest{Chart: model.ChartTypeDay, Fleet: true, Fans: true, Funds: true}

	res, err := makeRequest(rsiPledgeDataURL, data)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Pledge Data Pull Failed! Return Code %d|%s", res.StatusCode, res.Status)
	}

	// Ensure you have read the response body into a string variable before decoding it
	var responseBodyString string
	if res.Body != nil {
		defer res.Body.Close()
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		responseBodyString = string(bodyBytes)
	}

	payload := &model.RSIResponse{}
	err = json.Unmarshal([]byte(responseBodyString), &payload)
	if err != nil {
		return nil, err
	}

	if payload.Success == 0 {
		log.Printf("Raw JSON: %s\n", responseBodyString)
		return nil, fmt.Errorf("Pledge Data Pull Returned Unsuccessfully!\n\tCode: %v\n\tMessage: %v\n\tSuccess: %v\n", payload.Code, payload.Msg, payload.Success)
	}

	return &model.PledgeData{
		TimeStamp: timestamp,
		Funds:     payload.Data.Funds / 100,
		Citizens:  payload.Data.Fans,
		Fleet:     0, // not tracked anymore
	}, nil
}

func sendEmail(body string) {
	// Implement your email sending logic here
	// You can use a third-party email library or SMTP client to send the email
}

func main() {
	now := time.Now().UTC()
	min, hour, mday, mon, year := now.Minute(), now.Hour(), now.Day(), int(now.Month()), now.Year()
	currTime := fmt.Sprintf("%02d:%02d", hour, min)

	if min > 55 {
		hour++
	}

	timestamp := fmt.Sprintf("%04d-%02d-%02d %02d:00:00", year, mon, mday, hour)

	fmt.Printf("[%s GMT] RSI PledgeTracker Pull Running!\n", timestamp)

	doData := true
	if min > 5 && min < 55 {
		doData = false
		fmt.Printf("Tried to run at %s: is not within 5 minutes of the hour! Data insert disabled\n", currTime)
	}
	doData = true

	pledgeData, err := getPledgeData(timestamp)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		sendEmail(fmt.Sprintf("ERROR: %v\n", err))
		os.Exit(1)
	}

	if doData {
		db, err := database.ConnectPSQLFromEnv()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		pds := &service.PledgeDataService{DB: db}
		err = pds.Insert(pledgeData)
		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Printf("[%s GMT] Got %s: Funds %d Citizens %d Fleet %d\n", timestamp, timestamp, int(pledgeData.Funds), pledgeData.Citizens, pledgeData.Fleet)
}
