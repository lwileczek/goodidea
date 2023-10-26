package main

import (
	"log"
	"net/http"
)

func main() {
	err := connect()
	if err != nil {
		log.Fatal("Unable to create a connection to the database", err)
	}
	defer db.Close()

	mux := NewServer()
	log.Printf("Starting on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
