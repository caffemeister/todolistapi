package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TaskResponse struct {
	Message string `json:"message"`
	Task    Task   `json:"task"`
}

var (
	tasks     = make(map[int]Task)
	taskID    = 1
	tasksLock = &sync.Mutex{}
)

func nextID() int {
	taskID++
	return taskID - 1
}

// Use GET requests to get a list of tasks with params for ?status, ?title, ?limit (how many to show per page) and ?page (which page to go to)
// Use POST requests to create tasks using -d '{"Title": <title>}'
// Use PUT requests to update tasks using -d '{"Title": <newTitle>, "Completed": true}'
// Use DELETE requests to delete tasks by just specifying .../tasks/<id>

func main() {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		fmt.Println("")
		fmt.Printf("Completed in: %dms", elapsed.Milliseconds())
		fmt.Println("")
	}()
	defer db.Close()

	err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	http.HandleFunc("/tasks", TasksHandler)
	http.HandleFunc("/tasks/", TaskHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
