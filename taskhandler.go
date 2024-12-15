package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Handles specific tasks (by Task.id)
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	tasksLock.Lock()
	defer tasksLock.Unlock()

	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	fmt.Println(id)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	query := "SELECT id, title, completed FROM tasks WHERE id = ?"
	row := db.QueryRow(query, id)

	var task Task
	err = row.Scan(&task.ID, &task.Title, &task.Completed)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch task", http.StatusInternalServerError)
		return
	}

	switch r.Method {
	// PUT updates task with new entry
	case http.MethodPut:
		var updatedTask Task
		if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if updatedTask.Title != "" {
			task.Title = updatedTask.Title
		}
		if updatedTask.Completed != task.Completed {
			task.Completed = updatedTask.Completed
		}

		updateQuery := "UPDATE tasks SET title = ?, completed = ? WHERE id = ?"
		_, err := db.Exec(updateQuery, task.Title, task.Completed, task.ID)
		if err != nil {
			http.Error(w, "Failed to update task", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(task); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	// Deletes
	case http.MethodDelete:
		delete(tasks, id)

		w.WriteHeader(http.StatusNoContent) // No Content (204)
		fmt.Fprintf(w, "Task %d deleted successfully", id)

	default:
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
	}
}
