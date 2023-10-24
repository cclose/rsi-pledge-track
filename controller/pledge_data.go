package controller

import (
	"encoding/csv"
	"fmt"
	echo "github.com/labstack/echo/v4"
	logger "github.com/sirupsen/logrus"
	"github.org/cclose/rsi-pledge-track/model"
	"github.org/cclose/rsi-pledge-track/service"
	"net/http"
	"time"
)

var zeroTime = time.Time{}

type PledgeData struct {
	pledgeDataService service.IPledgeDataService
}

func NewPledgeDataController(pds service.IPledgeDataService) *PledgeData {
	return &PledgeData{
		pledgeDataService: pds,
	}
}

// GetPledgeData
// @Summary Get Pledge Data
// @Description Get Pledge Data.
// @Produce json html csv
// @Param timestamp query string false "Timestamp of a specific entry in RFC3339 format"
// @Param startingdatetime query string false "Get entries after this DateTime in RFC3339 format"
// @Param format query string false "Response format: json, html, csv" default("html")
// @Param offset query integer false "Timezone Offset from UTC" default(0)
// @Param limit query integer false "Limit for the number of results. Sending 0 will result in the default" default(100)
// @Success 200 {object} model.PledgeData "Successful response"
// @Failure 400 {string} string "Bad Request"
// @Failure 406 {string} string "Not Acceptable"
// @Failure 500 {string} string "Internal Server Error"
// @Router /pledge-data [get]
func (pdc *PledgeData) GetPledgeData(c echo.Context) error {
	gpdReq := &model.PledgeDataRequest{}
	err := gpdReq.ParseRequest(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var data []*model.PledgeData
	if gpdReq.TimeStamp != zeroTime {
		var dataEntry *model.PledgeData
		dataEntry, err = pdc.pledgeDataService.GetByTimestamp(gpdReq.TimeStamp, gpdReq.Offset)
		data = append(data, dataEntry)
	} else if gpdReq.AfterTimestamp != zeroTime {
		data, err = pdc.pledgeDataService.GetAfterTimestamp(gpdReq.AfterTimestamp, gpdReq.Offset, gpdReq.Limit)
	} else {
		data, err = pdc.pledgeDataService.GetAll(gpdReq.Limit, gpdReq.Offset)
	}
	if err != nil {
		logger.Error(err)
	}

	if gpdReq.Format == "json" || gpdReq.Format == "application/json" {
		return c.JSON(http.StatusOK, data)
	} else if gpdReq.Format == "csv" || gpdReq.Format == "text/csv" {
		// Implement CSV response
		return writeCSV(c, data)
	} else {
		// Implement HTML response
		htmlResponse := "<html>...</html>"
		return c.HTML(http.StatusOK, htmlResponse)
	}
}

func writeCSV(c echo.Context, data []*model.PledgeData) error {
	c.Response().Status = http.StatusOK
	c.Response().Header().Set("Content-Type", "text/csv")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=pledgeData.csv")

	writer := c.Response().Writer

	// Initialize a CSV writer
	csvWriter := csv.NewWriter(writer)

	// Write CSV header
	header := []string{"ID", "TimeStamp", "Funds", "Citizens", "Fleet"}
	_ = csvWriter.Write(header)
	csvWriter.Flush()

	for _, item := range data {
		// Write the data to the CSV writer
		_ = csvWriter.Write([]string{
			fmt.Sprintf("%d", item.ID),
			item.TimeStamp,
			fmt.Sprintf("%d", item.Funds),
			fmt.Sprintf("%d", item.Citizens),
			fmt.Sprintf("%d", item.Fleet),
		})
		csvWriter.Flush()
	}

	return nil
}
