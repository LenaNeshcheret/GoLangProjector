package main

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var (
	db  *sql.DB
	rdb *redis.Client
	mr  *miniredis.Miniredis
	ctx = context.Background()
)

// Setup initializes the in-memory database and Redis.
func setup() {
	var err error

	// Setup SQLite in-memory database
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create the table schema
	_, err = db.Exec(`
		CREATE TABLE tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			is_done BOOLEAN NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Setup miniredis
	mr, err = miniredis.Run()
	if err != nil {
		log.Fatalf("Failed to start miniredis: %v", err)
	}

	// Setup Redis client with miniredis address
	rdb = redis.NewClient(&redis.Options{
		Addr: mr.Addr(), // miniredis address
	})
}

// Teardown cleans up the database and Redis.
func teardown() {
	// Cleanup SQLite database
	db.Exec("DROP TABLE IF EXISTS tasks")

	// Cleanup miniredis
	mr.Close()
}

// TestGetAllTasksNoCache tests the scenario where data is not cached in Redis and is fetched from the database.
func TestGetAllTasksNoCache(t *testing.T) {
	setup()
	defer teardown()

	// Seed data into SQLite
	_, err := db.Exec("INSERT INTO tasks (title, is_done) VALUES (?, ?), (?, ?)", "Task 1", false, "Task 2", true)
	if err != nil {
		t.Fatalf("Failed to seed database: %v", err)
	}

	// Initialize TaskResource with the SQLite and Redis clients
	storage := Storage{
		db:    db,
		redis: rdb,
	}
	taskResource := TaskResource{
		s: &storage,
	}

	// Delete all keys from Redis to simulate cache miss
	rdb.FlushAll(ctx)

	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(taskResource.GetAll)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Expected response
	expected := `[{"id":1,"title":"Task 1","is_done":false},{"id":2,"title":"Task 2","is_done":true}]`

	// Check the response body
	assert.JSONEq(t, expected, rr.Body.String(), "handler returned unexpected body")

	// Ensure data is cached in Redis after the first request
	cachedData, err := rdb.Get(ctx, "all_tasks").Result()
	if err != nil {
		t.Fatalf("Failed to get data from Redis: %v", err)
	}
	if cachedData == "" {
		t.Errorf("Expected data to be cached in Redis")
	}
	//expectedCachedData := `[{"id":1,"title":"Task 1","is_done":false},{"id":2,"title":"Task 2","is_done":true}]`
	assert.JSONEq(t, expected, cachedData, "cache stores unexpected body")
}

func TestGetAllTasksCacheHit(t *testing.T) {
	setup()
	defer teardown()

	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Seed data into Redis to simulate cache hit
	rdb.Set(ctx, "all_tasks", `[{"id":1,"title":"Task 1","is_done":false},{"id":2,"title":"Task 2","is_done":true}]`, time.Hour)

	storage := Storage{
		db:    db,
		redis: rdb,
	}
	taskResource := TaskResource{
		s: &storage,
	}

	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(taskResource.GetAll)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Expected response
	expected := `[{"id":1,"title":"Task 1","is_done":false},{"id":2,"title":"Task 2","is_done":true}]`

	// Check the response body
	assert.JSONEq(t, expected, rr.Body.String(), "handler returned unexpected body")

	// Ensure that no queries were made to the database
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expected no database queries, but got some: %v", err)
	}
}
