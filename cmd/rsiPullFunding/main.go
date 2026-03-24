package main

import (
	"context"
	"github.org/cclose/rsi-pledge-track/controller"
	"log"
)

func main() {
	ctx := context.Background()
	res, err := controller.UpdatePledgeData(ctx)
	if res != "" {
		log.Println(res)
	}
	if err != nil {
		log.Fatal(err)
	}
}
