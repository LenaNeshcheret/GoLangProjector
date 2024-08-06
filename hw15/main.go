package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	connString := os.Getenv("POSTGRES_CONN_STR")
	if connString == "" {
		log.Fatal("POSTGRES_CONN_STR environment variable not set")
	}

	storage, err := NewStorage(connString)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Example usage
	tasks, err := storage.GetAllTasks()
	if err != nil {
		log.Fatalf("Failed to get all tasks: %v", err)
	}

	for _, task := range tasks {
		fmt.Printf("Task ID: %d, Title: %s, Done: %v\n", task.ID, task.Title, task.IsDone)
	}

	// Create an HTTP server and routes
	mux := http.NewServeMux()
	tasksResource := TaskResource{storage}
	mux.HandleFunc("GET /tasks", tasksResource.GetAll)
	mux.HandleFunc("POST /tasks", tasksResource.CreateOne)
	mux.HandleFunc("PUT /tasks/{id}", tasksResource.UpdateOne)
	mux.HandleFunc("DELETE /tasks/{id}", tasksResource.DeleteOne)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	IsDone bool   `json:"is_done"`
}

type Storage struct {
	db *sql.DB
}

func NewStorage(connString string) (*Storage, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error during pinging database: %w", err)
	}

	_, err = db.Query("CREATE table IF NOT EXISTS  tasks (id serial, title varchar, is_done boolean)")
	if err != nil {
		return nil, fmt.Errorf("error during during taable creation: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetAllTasks() ([]Task, error) {
	rows, err := s.db.Query("SELECT id, title, is_done FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("selecting tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.IsDone)
		if err != nil {
			return nil, fmt.Errorf("scanning rows: %w", err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return tasks, nil
}

func (s *Storage) CreateOneTask(task Task) (int, error) {
	var id int
	err := s.db.QueryRow("INSERT INTO tasks (title, is_done) VALUES ($1, $2) RETURNING id", task.Title, task.IsDone).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting task: %w", err)
	}
	return id, nil
}

func (s *Storage) UpdateTask(id int, task Task) (bool, error) {
	result, err := s.db.Exec("UPDATE tasks SET title = $1, is_done = $2 WHERE id = $3", task.Title, task.IsDone, id)
	if err != nil {
		return false, fmt.Errorf("updating task: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("checking rows affected: %w", err)
	}
	return rowsAffected > 0, nil
}

func (s *Storage) DeleteTaskByID(id int) (bool, error) {
	result, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return false, fmt.Errorf("deleting task: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("checking rows affected: %w", err)
	}
	return rowsAffected > 0, nil
}

type TaskResource struct {
	s *Storage
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
		http.Error(w, "Failed to decode task", http.StatusBadRequest)
		return
	}

	id, err := t.s.CreateOneTask(task)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}
	task.ID = id

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, "Failed to encode task", http.StatusInternalServerError)
	}
}

func (t *TaskResource) UpdateOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var updatedTask Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Failed to decode task", http.StatusBadRequest)
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
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(updatedTask)
	if err != nil {
		http.Error(w, "Failed to encode task", http.StatusInternalServerError)
	}
}

func (t *TaskResource) DeleteOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idVal := strings.TrimPrefix(r.URL.Path, "/tasks/")
	taskID, err := strconv.Atoi(idVal)
	if err != nil {
		http.Error(w, "Invalid id param", http.StatusBadRequest)
		return
	}

	ok, err := t.s.DeleteTaskByID(taskID)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
}
