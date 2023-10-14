package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.org/cclose/rsi-pledge-track/controller"
	"github.org/cclose/rsi-pledge-track/database"
	"github.org/cclose/rsi-pledge-track/service"
	"os"
	"strconv"
	"time"
)

var logger *logrus.Logger

func init() {
	// Create a new Logrus logger instance
	logger = logrus.New()

	// Set the log level (e.g., Info, Debug, Error, etc.)
	logger.SetLevel(logrus.InfoLevel)

	// Set the log format (text or JSON)
	logger.SetFormatter(&logrus.TextFormatter{})

	// Set the output destination (os.Stdout for console, os.Stderr for stderr)
	logger.SetOutput(os.Stdout)
}

func main() {
	if startDelay := os.Getenv("START_DELAY"); startDelay != "" {
		if delay, err := strconv.Atoi(startDelay); err == nil {
			logger.Infof("START_DELAY %s... sleeping", startDelay)
			time.Sleep(time.Duration(delay) * time.Second)
			logger.Info("... resuming")
		} else {
			logger.Warnf("Invalids START_DELAY %s", startDelay)
		}
	}

	db, err := database.ConnectPSQLFromEnv()
	if err != nil {
		logger.Fatal(err)
	}

	pds := service.NewPledgeDataService(db)
	pdc := controller.NewPledgeDataController(pds)

	echoPort := os.Getenv("PORT")
	if echoPort == "" {
		echoPort = "8080"
	}

	e := echo.New()

	// Route to retrieve data in different formats
	e.GET("/pledge-data", pdc.GetPledgeData)

	//go func() {
	if err := e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		logger.Info("shutting down the server")
	}
	//}()
}
