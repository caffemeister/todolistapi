package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Handles the showing of tasks and the adding of tasks
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	// GET shows the list of current tasks
	if r.Method == http.MethodGet {
		statusFilter := r.URL.Query().Get("status") // gets the value of tasks?status=<value>
		titleFilter := r.URL.Query().Get("title")
		limitManual := r.URL.Query().Get("limit")
		pageManual := r.URL.Query().Get("page")

		// Default values for page and limit per page
		limitDefault := 10
		pageDefault := 1

		// This checks if a limit-per-page value was provided
		if limitManual != "" {
			if limitCheck, err := strconv.Atoi(limitManual); err == nil && limitCheck > 0 {
				limitDefault = limitCheck
			}
		}

		// This checks if a page value was provided
		if pageManual != "" {
			if pageCheck, err := strconv.Atoi(pageManual); err == nil && pageCheck > 0 {
				pageDefault = pageCheck
			}
		}

		offset := (pageDefault - 1) * limitDefault

		// 1=1 just means select everything
		query := "SELECT id, title, completed FROM tasks WHERE 1=1"
		args := []interface{}{}

		if statusFilter != "" {
			switch statusFilter {
			case "completed":
				query += " AND completed = ?"
				args = append(args, true)
			case "pending":
				query += " AND completed = ?"
				args = append(args, false)
			default:
			}
		}

		if titleFilter != "" {
			query += " AND LOWER(title) LIKE ?"
			args = append(args, "%"+strings.ToLower(titleFilter)+"%")
		}

		// Pagination here
		query += " LIMIT ? OFFSET ?"
		args = append(args, limitDefault, offset)

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var task Task
			if err := rows.Scan(&task.ID, &task.Title, &task.Completed); err != nil {
				http.Error(w, "Failed to parse tasks", http.StatusInternalServerError)
				return
			}
			tasks = append(tasks, task)
		}

		// Send response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// POST adds a new task to current tasks
	} else if r.Method == http.MethodPost {
		var newTask Task
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		newTask.ID = nextID()
		tasks[newTask.ID] = newTask

		response := TaskResponse{
			Message: "Created a new task",
			Task:    newTask,
		}

		insertString := `INSERT INTO tasks (title, completed) VALUES ($1, $2)`
		_, err := db.Exec(insertString, newTask.Title, newTask.Completed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}
