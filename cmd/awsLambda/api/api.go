package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	logger "github.com/sirupsen/logrus"
	"github.org/cclose/rsi-pledge-track/controller"
	"github.org/cclose/rsi-pledge-track/database"
	"github.org/cclose/rsi-pledge-track/service"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := database.ConnectPSQLFromEnv()
	if err != nil {
		logger.Fatal(err)
	}

	pds := service.NewPledgeDataService(db)
	pdc := controller.NewPledgeDataController(pds)
	e.Renderer = pdc

	e.GET("/", pdc.GetPledgeData)

	echoLambda = echoadapter.New(e)
}

func HandleAPIRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, request)
}

func main() {
	lambda.Start(HandleAPIRequest)
}
