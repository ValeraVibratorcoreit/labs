package main

import "time"

type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

type createTaskRequest struct {
	Title string `json:"title"`
}

type updateTaskRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}
