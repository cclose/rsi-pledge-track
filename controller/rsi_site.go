package controller

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.org/cclose/rsi-pledge-track/database"
	"github.org/cclose/rsi-pledge-track/model"
	"github.org/cclose/rsi-pledge-track/service"
	"io/ioutil"
	"log"
	"net/http"
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

func PullRSIPledgeData(timestamp string) (*model.PledgeData, error) {
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
		Funds:     payload.Data.Funds,
		Citizens:  int32(payload.Data.Fans),
		Fleet:     0, // not tracked anymore
	}, nil
}

func CheckEventTime(hour, minute int) error {
	if minute > 5 && minute < 55 {
		return fmt.Errorf("Tried to run at %02d:%02d: is not within 5 minutes of the hour!\n",
			hour, minute)
	}

	return nil
}

func PullPledgeData() (*model.PledgeData, error) {
	now := time.Now().UTC()
	minute, hour, mday, mon, year := now.Minute(), now.Hour(), now.Day(), int(now.Month()), now.Year()

	// Check if we need to adjust hour
	if minute > 55 {
		hour++
	}

	// Format timestamp
	timestamp := fmt.Sprintf("%04d-%02d-%02d %02d:00:00", year, mon, mday, hour)
	log.Printf("[%s GMT] RSI PledgeTracker Pull Running!\n", timestamp)

	// Pull data
	pledgeData, err := PullRSIPledgeData(timestamp)
	if err != nil {
		return nil, fmt.Errorf("[PullRSIPledgeData] error: %w", err)
	}

	return pledgeData, CheckEventTime(hour, minute)
}

func PushPledgeData(ctx context.Context, pledgeData *model.PledgeData) error {
	// Connect to database
	db, err := database.ConnectPSQLFromEnv()
	if err != nil {
		return fmt.Errorf("[ConnectPSQL] error: %w", err)
	}
	defer db.Close()

	// Parse the timestamp string into a time.Time object
	timestamp, err := time.Parse("2006-01-02 15:04:05", pledgeData.TimeStamp)
	if err != nil {
		return fmt.Errorf("Error parsing timestamp: %v", err)
	}

	// Extract hour and minute from the timestamp  and verify it's good
	if err = CheckEventTime(timestamp.Hour(), timestamp.Minute()); err != nil {
		return err
	}

	// Insert data into database
	pds := &service.PledgeDataService{DB: db}
	err = pds.Insert(ctx, pledgeData)
	if err != nil {
		return fmt.Errorf("[PledgeDataService.Insert] error: %w", err)
	}

	return nil
}

func UpdatePledgeData(ctx context.Context) (string, error) {
	// Pull data
	pledgeData, err := PullPledgeData()
	if err != nil {
		return "", err
	}

	// Push data
	err = PushPledgeData(ctx, pledgeData)
	if err != nil {
		return "", err
	}

	// Format result
	result := fmt.Sprintf("[%s GMT] Got %s: Funds %d Citizens %d Fleet %d\n", pledgeData.TimeStamp,
		pledgeData.TimeStamp, pledgeData.Funds, pledgeData.Citizens, pledgeData.Fleet)
	log.Printf("Pulled %s\n", result)

	return result, nil
}
