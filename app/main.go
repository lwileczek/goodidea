package main

import (
	"log"
	"net/http"

	"github.com/lwileczek/goodidea"
)

func main() {
	err := goodidea.Connect()
	if err != nil {
		log.Fatal("Unable to create a connection to the database", err)
	}
	defer goodidea.DB.Close()

	//Create a logger
	goodidea.SetupLogger()

	//Set up mux router
	mux := goodidea.NewServer()

	//Runnit
	log.Printf("Starting on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
