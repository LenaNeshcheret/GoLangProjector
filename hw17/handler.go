package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type taskStorage interface {
	GetAllTasks() ([]Task, error)
	GetTaskByID(id int) (Task, error)
	CreateTask(task Task) (int, error)
	UpdateTask(id int, task Task) (bool, error)
	DeleteTaskByID(id int) (bool, error)
}

type TaskResource struct {
	s taskStorage
}

func (t *TaskResource) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tasks, err := t.s.GetAllTasks()
	if err != nil {
		http.Error(w, "Failed to get tasks", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
	}
}

func (t *TaskResource) CreateOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Failed to decode task, error:"+err.Error(), http.StatusBadRequest)
		return
	}

	id, err := t.s.CreateTask(task)
	if err != nil {
		http.Error(w, "Failed to create task, error:"+err.Error(), http.StatusInternalServerError)
		return
	}
	task.ID = id

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, "Failed to encode task, error:"+err.Error(), http.StatusInternalServerError)
	}
}

func (t *TaskResource) UpdateById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var updatedTask Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Failed to decode task, error:"+err.Error(), http.StatusBadRequest)
		return
	}

	idVal := strings.TrimPrefix(r.URL.Path, "/tasks/")
	taskID, err := strconv.Atoi(idVal)
	if err != nil {
		http.Error(w, "Invalid id param", http.StatusBadRequest)
		return
	}

	ok, err := t.s.UpdateTask(taskID, updatedTask)
	if err != nil {
		http.Error(w, "Failed to update task, error:"+err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(updatedTask)
	if err != nil {
		http.Error(w, "Failed to encode task, error:"+err.Error(), http.StatusInternalServerError)
	}
}

func (t *TaskResource) DeleteByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idVal := strings.TrimPrefix(r.URL.Path, "/tasks/")
	taskID, err := strconv.Atoi(idVal)
	if err != nil {
		http.Error(w, "Invalid id param, error:"+err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := t.s.DeleteTaskByID(taskID)
	if err != nil {
		http.Error(w, "Failed to delete task, error:"+err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
}

func (t *TaskResource) GetById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idVal := strings.TrimPrefix(r.URL.Path, "/tasks/")
	taskID, err := strconv.Atoi(idVal)
	if err != nil {
		http.Error(w, "Invalid id param", http.StatusBadRequest)
		return
	}

	task, err := t.s.GetTaskByID(taskID)
	if err != nil {
		http.Error(w, "Failed to get task with id "+string(int32(taskID))+", error:"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the task is empty, meaning not found
	if task.ID == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Failed to encode task to JSON", http.StatusInternalServerError)
		return
	}
}
