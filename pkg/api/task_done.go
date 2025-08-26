package api

import (
	"net/http"
	"time"

	"LAST_TODO_2/pkg/db"
)

// doneHandler обрабатывает отметку задачи выполненной
func doneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		return
	}

	// одноразовая задача — удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{}, http.StatusOK)
		return
	}

	// для тестов можно передавать now
	nowStr := r.URL.Query().Get("now")
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateLayout, nowStr)
		if err != nil {
			writeJSON(w, map[string]string{"error": "invalid now format"}, http.StatusBadRequest)
			return
		}
	}

	// вызываем функцию из api/nextdate.go
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}

// deleteHandler обрабатывает удаление задачи
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{}, http.StatusOK)
}
