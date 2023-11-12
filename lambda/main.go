package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/lwileczek/goodidea"
)

var app *gorillamux.GorillaMuxAdapterV2

func init() {
	if goodidea.DB == nil {
		err := goodidea.Connect()
		if err != nil {
			log.Fatal("Unable to create a connection to the database", err)
		}
	}
	//defer goodidea.DB.Close()

	if goodidea.Logr == nil {
		goodidea.SetupLogger()
	}

	//Set up mux router
	router := goodidea.NewServer()

	app = gorillamux.NewV2(router)
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return app.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
