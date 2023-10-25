package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	mux := NewServer()
	log.Printf("Starting on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

type Task struct {
	ID          uint32    `json:"id"`
	Status      bool      `json:"status"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Score       int32     `json:"score"`
	CompletedAt time.Time `json:"completedAt"`
	CreatedAt   time.Time `json:"createdAt"`
	DeletedAt   time.Time `json:"deletedAt"`
}
