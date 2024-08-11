package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// MockStorage simulates the storage layer for testing.
type MockStorage struct {
	Storage
	mockCtrl        *gomock.Controller
	mockGetAllTasks func() ([]Task, error)
}

// GetAllTasks returns tasks or an error based on the test case setup.
func (m *MockStorage) GetAllTasks() ([]Task, error) {
	return m.mockGetAllTasks()
}

// TestGetAllTasksSuccess tests the scenario where tasks are successfully retrieved.
func TestGetAllTasksSuccess(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := &MockStorage{
		mockCtrl: mockCtrl,
		mockGetAllTasks: func() ([]Task, error) {
			return []Task{
				{ID: 1, Title: "Test Task 1", IsDone: false},
				{ID: 2, Title: "Test Task 2", IsDone: true},
			}, nil
		},
	}

	// Initialize TaskResource with the mock storage
	taskResource := TaskResource{s: mockStorage}

	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(taskResource.GetAll)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := []Task{
		{ID: 1, Title: "Test Task 1", IsDone: false},
		{ID: 2, Title: "Test Task 2", IsDone: true},
	}

	var actual []Task
	err = json.NewDecoder(rr.Body).Decode(&actual)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(actual) != len(expected) {
		t.Errorf("expected %d tasks, got %d", len(expected), len(actual))
	}
}

// TestGetAllTasksError tests the scenario where an error occurs while retrieving tasks.
func TestGetAllTasksError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := &MockStorage{
		mockCtrl: mockCtrl,
		mockGetAllTasks: func() ([]Task, error) {
			return nil, fmt.Errorf("error fetching tasks")
		},
	}

	// Initialize TaskResource with the mock storage
	taskResource := TaskResource{s: mockStorage}

	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(taskResource.GetAll)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := "Failed to get tasks"
	assert.Equal(t, expected, strings.TrimSpace(rr.Body.String()), "handler returned unexpected body")
}

func TestInvalidMethod(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorage := &MockStorage{
		mockCtrl: mockCtrl,
	}

	taskResource := TaskResource{s: mockStorage}

	// List of HTTP methods to test
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		req, err := http.NewRequest(method, "/tasks", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(taskResource.GetAll)
		handler.ServeHTTP(rr, req)

		// Assert status code
		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "handler returned incorrect status code for method "+method)

		// Assert response body (optional)
		expectedBody := ""
		assert.Equal(t, expectedBody, strings.TrimSpace(rr.Body.String()), "handler returned unexpected body for method "+method)
	}
}
