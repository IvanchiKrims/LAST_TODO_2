package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"LAST_TODO_2/pkg/db"
)

// addTaskHandler обрабатывает POST /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()}, http.StatusBadRequest) //  Обработка кода ответа
		return
	}

	// Проверка обязательного заголовка
	if task.Title == "" {
		writeJSON(w, map[string]any{"error": "Не указан заголовок задачи"}, http.StatusBadRequest) //  обработка кода ответа
		return
	}

	// Проверка и корректировка даты
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]any{"error": err.Error()}, http.StatusBadRequest) //  обработка кода ответа
		return
	}

	// Добавляем в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error()}, http.StatusInternalServerError) // обработка кода ответа
		return
	}

	writeJSON(w, map[string]any{"id": fmt.Sprint(id)}, http.StatusOK) //  обработка кода ответа
}

// checkDate проверяет и корректирует дату задачи
func checkDate(task *db.Task) error {
	now := time.Now()

	// Если дата не указана — ставим сегодня
	if task.Date == "" {
		task.Date = now.Format(DateLayout)
	}

	// Проверяем формат даты
	t, err := time.Parse(DateLayout, task.Date)
	if err != nil {
		return err
	}

	// Проверка правила повторения (если есть)
	if task.Repeat != "" {
		_, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// Если дата уже прошла
	if t.Before(truncateToDay(now)) {
		if task.Repeat == "" {
			task.Date = now.Format(DateLayout)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}
	return nil
}

// truncateToDay — обрезает время до даты
func truncateToDay(t time.Time) time.Time {
	res, _ := time.Parse(DateLayout, t.Format(DateLayout))
	return res
}

// writeJSON — утилита для возврата JSON с кодом ответа и обработкой ошибки кодирования
func writeJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)                                   //  установка кода ответа
	if err := json.NewEncoder(w).Encode(data); err != nil { //  обработка ошибки JSON
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
