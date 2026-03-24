package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.org/cclose/rsi-pledge-track/controller"
)

func HandlePullRequest(ctx context.Context) (*string, error) {
	res, err := controller.PullPledgeData()
	resJSON := ""
	if res != nil {
		resJSONBytes, jerr := json.Marshal(res)
		if jerr != nil {
			err = fmt.Errorf("Multiple Errors while Marshalling: err %w, json_err: %w", err, jerr)
		} else {
			resJSON = string(resJSONBytes)
		}
	}
	return &resJSON, err
}

func main() {
	lambda.Start(HandlePullRequest)
}
