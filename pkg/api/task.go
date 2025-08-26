package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"LAST_TODO_2/pkg/db"
)

// taskHandler обрабатывает GET, POST, PUT, DELETE для /api/task
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			writeJSON(w, map[string]any{"error": "Не указан идентификатор"}, http.StatusBadRequest)
			return
		}

		task, err := db.GetTask(id)
		if err != nil {
			writeJSON(w, map[string]any{"error": "Задача не найдена"}, http.StatusNotFound)
			return
		}

		resp := map[string]any{
			"id":      fmt.Sprint(task.ID),
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}
		writeJSON(w, resp, http.StatusOK)

	case http.MethodPut:
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeJSON(w, map[string]any{"error": err.Error()}, http.StatusBadRequest)
			return
		}

		rawID, ok := payload["id"]
		if !ok {
			writeJSON(w, map[string]any{"error": "Не указан идентификатор"}, http.StatusBadRequest)
			return
		}

		idStr := fmt.Sprint(rawID)
		taskID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || taskID <= 0 {
			writeJSON(w, map[string]any{"error": "Некорректный идентификатор"}, http.StatusBadRequest)
			return
		}

		task := db.Task{
			ID:      taskID,
			Date:    fmt.Sprint(payload["date"]),
			Title:   fmt.Sprint(payload["title"]),
			Comment: fmt.Sprint(payload["comment"]),
			Repeat:  fmt.Sprint(payload["repeat"]),
		}

		if task.Title == "" {
			writeJSON(w, map[string]any{"error": "Не указан заголовок"}, http.StatusBadRequest)
			return
		}

		if err := checkDate(&task); err != nil {
			writeJSON(w, map[string]any{"error": err.Error()}, http.StatusBadRequest)
			return
		}

		if err := db.UpdateTask(&task); err != nil {
			writeJSON(w, map[string]any{"error": "Задача не найдена"}, http.StatusNotFound)
			return
		}

		writeJSON(w, map[string]any{}, http.StatusOK)

	case http.MethodDelete:
		deleteHandler(w, r)

	default:
		writeJSON(w, map[string]any{"error": "Неподдерживаемый метод"}, http.StatusMethodNotAllowed)
	}
}
