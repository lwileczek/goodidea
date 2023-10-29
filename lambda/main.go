package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lwileczek/goodidea"
)

func main() {
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

	app := httpadapter.New(handlers.LoggingHandler(os.Stdout, router))

	lambda.Start(app.ProxyWithContext)
}
