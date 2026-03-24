package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.org/cclose/rsi-pledge-track/controller"
	"github.org/cclose/rsi-pledge-track/model"
	"log"
	"time"
)

// LambdaInvocationRecord represents the root object of the JSON.
type LambdaInvocationRecord struct {
	Version         string                 `json:"version"`
	Timestamp       time.Time              `json:"timestamp"`
	RequestContext  RequestContext         `json:"requestContext"`
	RequestPayload  map[string]interface{} `json:"requestPayload"`
	ResponseContext ResponseContext        `json:"responseContext"`
	ResponsePayload string                 `json:"responsePayload"`
}

// RequestContext represents the request context object within the JSON.
type RequestContext struct {
	RequestID              string `json:"requestId"`
	FunctionARN            string `json:"functionArn"`
	Condition              string `json:"condition"`
	ApproximateInvokeCount int    `json:"approximateInvokeCount"`
}

// ResponseContext represents the response context object within the JSON.
type ResponseContext struct {
	StatusCode      int    `json:"statusCode"`
	ExecutedVersion string `json:"executedVersion"`
}

func HandlePushRequest(ctx context.Context, event json.RawMessage) (*string, error) {
	// Print the received invocation record
	log.Println("Received invocation record")
	var lambdaRecord LambdaInvocationRecord
	err := json.Unmarshal(event, &lambdaRecord)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall request event: %w", err)
	}

	// Decode the payload from input data
	reqJSON := lambdaRecord.ResponsePayload
	log.Printf("Received: [%s]\n", reqJSON)
	var pledgeData model.PledgeData
	err = json.Unmarshal([]byte(reqJSON), &pledgeData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall request payload: %w", err)
	}

	// insert our data
	err = controller.PushPledgeData(ctx, &pledgeData)

	res := "Push executed"
	return &res, err
}

func main() {
	lambda.Start(HandlePushRequest)
}
