package controller

import (
	echo "github.com/labstack/echo/v4"
	logger "github.com/sirupsen/logrus"
	"github.org/cclose/rsi-pledge-track/service"
	"net/http"
)

type PledgeData struct {
	pledgeDataService *service.PledgeDataService
}

func NewPledgeDataController(pds *service.PledgeDataService) *PledgeData {
	return &PledgeData{
		pledgeDataService: pds,
	}
}
func (pdc *PledgeData) GetPledgeData(c echo.Context) error {
	// Determine the response format based on the Accept header or 'format' query parameter
	format := c.Request().Header.Get("Accept")
	if format == "" {
		// If Accept header is not set, check the 'format' query parameter
		format = c.QueryParam("format")
	}
	logger.Infof("format |%s|", format)

	data, err := pdc.pledgeDataService.GetAll()
	if err != nil {
		logger.Error(err)
	}

	if format == "json" || format == "application/json" {
		return c.JSON(http.StatusOK, data)
	} else if format == "csv" || format == "text/csv" {
		// Implement CSV response
		return c.String(http.StatusOK, "CSV data")
	} else {
		// Implement HTML response
		htmlResponse := "<html>...</html>"
		return c.HTML(http.StatusOK, htmlResponse)
	}
}
