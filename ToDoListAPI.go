package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"

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
	taskID    int
	tasksLock = &sync.Mutex{}
)

func nextID() int {
	var taskID int

	// Query the database for the MAX(id) in tasks
	err := db.QueryRow("SELECT MAX(id) FROM tasks").Scan(&taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found in the tasks table.")
		} else {
			panic(err) // Handle unexpected errors
		}
	}
	taskID++
	return taskID
}

// Use GET requests to get a list of tasks with params for ?status, ?title, ?limit (how many to show per page) and ?page (which page to go to)
// Use POST requests to create tasks using -d '{"Title": <title>}'
// Use PUT requests to update tasks using -d '{"Title": <newTitle>, "Completed": true}'
// Use DELETE requests to delete tasks by just specifying .../tasks/<id>

func main() {
	defer db.Close()

	err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	taskID = fetchID()
	http.HandleFunc("/tasks", TasksHandler)
	http.HandleFunc("/tasks/", TaskHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func fetchID() int {
	err := db.QueryRow("SELECT MAX(id) FROM tasks").Scan(&taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Can't find task IDs!")
		} else {
			panic(err)
		}
	}
	return taskID
}
