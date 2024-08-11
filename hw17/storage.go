package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	IsDone bool   `json:"is_done"`
}

type Storage struct {
	db    *sql.DB
	redis *redis.Client
}

func NewStorage(connPostgres string, connRedis string) (*Storage, error) {
	//postgres
	db, err := createPostgres(connPostgres)
	if err != nil {
		return nil, fmt.Errorf("error during create database postgres: %w", err)
	}
	//redis

	redisClient, err := createRedis(connRedis, err)
	if err != nil {
		return nil, fmt.Errorf("error during create database redis: %w", err)
	}
	return &Storage{
		db:    db,
		redis: redisClient,
	}, nil
}

func createRedis(connRedis string, err error) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: connRedis,
	})

	ctx := context.Background()

	res, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Pinging redis: %v", err)
	}
	log.Printf("Pinged: %v", res)
	return redisClient, err
}

func createPostgres(connPostgres string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connPostgres)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error during pinging database: %w", err)
	}
	return db, err
}

func (s *Storage) GetAllTasks() ([]Task, error) {
	ctx := context.Background()
	cacheKey := "all_tasks"

	// Try fetch the tasks from redis
	var tasks []Task
	found, err := s.getFromCache(ctx, cacheKey, &tasks)
	if found {
		if err == nil {
			return tasks, nil
		} else {
			return nil, err
		}
	}

	// Fetch from DB if tasks do not exist in redis
	tasks, err = s.fetchAllTasksFromDB()
	if err != nil {
		return nil, err
	}

	// Cache tasks to redis
	err = s.setToCache(ctx, cacheKey, tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *Storage) GetTaskByID(id int) (Task, error) {
	ctx := context.Background()
	cacheKey := "tasks_" + strconv.Itoa(id)

	// Try fetch the task from redis
	var task Task
	found, err := s.getFromCache(ctx, cacheKey, &task)
	if found {
		if err == nil {
			return task, nil
		} else {
			return Task{}, err
		}
	}

	// Fetch from DB if task does not exist in redis
	task, err = s.fetchTaskByIDFromDB(id)
	if err != nil {
		return Task{}, err
	}

	// Cache task to redis
	err = s.setToCache(ctx, cacheKey, task)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (s *Storage) CreateTask(task Task) (int, error) {
	var id int
	// fetch from db
	err := s.db.QueryRow("INSERT INTO tasks (title, is_done) VALUES ($1, $2) RETURNING id", task.Title, task.IsDone).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting task: %w", err)
	}
	// invalidate cache for all_tasks
	_, err = s.redis.Del(context.Background(), "all_tasks").Result()
	if err != nil {
		log.Printf("could not invalidate cache for all_tasks key")
	}
	return id, nil
}

func (s *Storage) UpdateTask(id int, task Task) (bool, error) {

	updated, err := s.updateTask(id, task)
	if err != nil {
		return updated, err
	}
	s.invalidateCache(id, updated)
	return updated, nil
}

func (s *Storage) DeleteTaskByID(id int) (bool, error) {
	isDeleted, err := s.deleteTask(id)
	if err != nil {
		return false, err
	}
	s.invalidateCache(id, isDeleted)
	return isDeleted, nil
}
func (s *Storage) getFromCache(ctx context.Context, cacheKey string, data interface{}) (bool, error) {
	cachedData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return true, fmt.Errorf("getting data from redis: %w", err)
	}

	err = json.Unmarshal([]byte(cachedData), data)
	if err != nil {
		return true, fmt.Errorf("unmarshaling data from redis: %w", err)
	}

	return true, nil
}

func (s *Storage) setToCache(ctx context.Context, cacheKey string, data interface{}) error {
	taskData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling data: %w", err)
	}
	err = s.redis.Set(ctx, cacheKey, taskData, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("setting data in redis: %w", err)
	}
	return nil
}

func (s *Storage) fetchAllTasksFromDB() ([]Task, error) {
	rows, err := s.db.Query("SELECT id, title, is_done FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("selecting tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Title, &task.IsDone)
		if err != nil {
			return nil, fmt.Errorf("scanning rows: %w", err)
		}
		tasks = append(tasks, task)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return tasks, nil
}

func (s *Storage) fetchTaskByIDFromDB(id int) (Task, error) {
	row := s.db.QueryRow("SELECT id, title, is_done FROM tasks WHERE id = $1", id)
	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.IsDone)
	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, nil
		}
		return Task{}, fmt.Errorf("scanning row: %w", err)
	}

	return task, nil
}

func (s *Storage) deleteTask(id int) (bool, error) {
	result, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return false, fmt.Errorf("deleting task: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("checking rows affected: %w", err)
	}
	isDeleted := rowsAffected > 0
	return isDeleted, nil
}

func (s *Storage) invalidateCache(id int, updated bool) {
	if updated {
		cacheKey := "tasks_" + strconv.Itoa(id)
		_, err := s.redis.Del(context.Background(), "all_tasks", cacheKey).Result()
		if err != nil {
			log.Printf("could not invalidate cache for all_tasks," + cacheKey)
		}
	}
}

func (s *Storage) updateTask(id int, task Task) (bool, error) {
	result, err := s.db.Exec("UPDATE tasks SET title = $1, is_done = $2 WHERE id = $3", task.Title, task.IsDone, id)
	if err != nil {
		return false, fmt.Errorf("updating task: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("checking rows affected: %w", err)
	}
	updated := rowsAffected > 0
	return updated, nil
}
