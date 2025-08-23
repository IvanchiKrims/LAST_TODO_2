package api

import (
	"net/http"
	"time"

	"LAST_TODO_2/pkg/db"
)

// doneHandler обрабатывает отметку задачи выполненной
func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]any{"error": "method not allowed"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]any{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]any{"error": "Задача не найдена"})
		return
	}

	// одноразовая задача — удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]any{"error": err.Error()})
			return
		}
		writeJSON(w, map[string]any{})
		return
	}

	// для тестов можно передавать now
	nowStr := r.URL.Query().Get("now")
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			writeJSON(w, map[string]any{"error": "invalid now format"})
			return
		}
	}

	// вызываем функцию из api/nextdate.go
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]any{})
}

// deleteHandler обрабатывает удаление задачи
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, map[string]any{"error": "method not allowed"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]any{"error": "Не указан идентификатор"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]any{})
}
