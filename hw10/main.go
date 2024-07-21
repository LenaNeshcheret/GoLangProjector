package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	mux := http.NewServeMux()

	s := NewStorage()

	tasks := TaskResource{
		s: s,
	}

	mux.HandleFunc("GET /tasks", tasks.GetAll)
	mux.HandleFunc("POST /tasks", tasks.CreateOne)
	mux.HandleFunc("PUT /tasks/{id}", tasks.UpdateOne)
	mux.HandleFunc("DELETE /tasks/{id}", tasks.DeleteOne)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Failed to listen and serve: %v\n", err)
	}
}

type TaskResource struct {
	s *Storage
}

func (t *TaskResource) GetAll(w http.ResponseWriter, r *http.Request) {
	trips := t.s.GetAllTasks()

	err := json.NewEncoder(w).Encode(trips)
	if err != nil {
		fmt.Printf("Failed to encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *TaskResource) CreateOne(w http.ResponseWriter, r *http.Request) {
	var task Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		fmt.Printf("Failed to decode: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task.ID = t.s.CreateOneTask(task)

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		fmt.Printf("Failed to encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *TaskResource) UpdateOne(w http.ResponseWriter, r *http.Request) {
	var updatedTask Task

	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		fmt.Printf("Failed to decode: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idVal := strings.TrimPrefix(r.URL.Path, "/tasks/")
	taskID, err := strconv.Atoi(idVal)
	if err != nil {
		fmt.Printf("Invalid id param: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok := t.s.UpdateTask(taskID, updatedTask)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(updatedTask)
	if err != nil {
		fmt.Printf("Failed to encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *TaskResource) DeleteOne(w http.ResponseWriter, r *http.Request) {
	idVal := r.PathValue("id")

	taskID, err := strconv.Atoi(idVal)
	if err != nil {
		fmt.Printf("Invalid id param: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok := t.s.DeleteTaskByID(taskID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
