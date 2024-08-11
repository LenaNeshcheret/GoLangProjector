package main

import (
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {

	//redis
	connRedis := os.Getenv("REDIS_CONN_STR")
	if connRedis == "" {
		log.Fatal("REDIS_CONN_STR environment variable not set")
	}

	// postgres
	//connPostgres := "user=postgres password=password dbname=university host=localhost port=5432 sslmode=disable"
	connPostgres := os.Getenv("POSTGRES_CONN_STR")
	if connPostgres == "" {
		log.Fatal("POSTGRES_CONN_STR environment variable not set")
	}

	storage, err := NewStorage(connPostgres, connRedis)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Create an HTTP server and routes
	mux := http.NewServeMux()

	tasksResource := TaskResource{storage}
	mux.HandleFunc("GET /tasks", tasksResource.GetAll)
	mux.HandleFunc("GET /tasks/{id}", tasksResource.GetById)
	mux.HandleFunc("POST /tasks", tasksResource.CreateOne)
	mux.HandleFunc("PUT /tasks/{id}", tasksResource.UpdateById)
	mux.HandleFunc("DELETE /tasks/{id}", tasksResource.DeleteByID)

	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}
