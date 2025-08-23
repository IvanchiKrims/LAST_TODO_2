package api

import (
	"fmt"
	"net/http"
	"time"

	"LAST_TODO_2/pkg/db"
)

// TasksResp — обёртка для JSON ответа
type TaskResp struct {
	ID      string `json:"id"` // ID конвертирован в строку
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// tasksHandler обрабатывает GET /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	var tasks []*db.Task
	var err error

	if search == "" {
		tasks, err = db.Tasks(50) // получаем максимум 50 ближайших задач
	} else {
		// Проверяем, может ли search быть датой в формате 02.01.2006
		if t, errDate := time.Parse("02.01.2006", search); errDate == nil {
			date := t.Format("20060102")
			tasks, err = db.TasksByDate(date, 50)
		} else {
			tasks, err = db.TasksBySearch(search, 50)
		}
	}

	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Если срез nil, создаём пустой
	if tasks == nil {
		tasks = []*db.Task{}
	}

	// Преобразуем задачи в формат со строковым ID
	respTasks := make([]TaskResp, len(tasks))
	for i, t := range tasks {
		respTasks[i] = TaskResp{
			ID:      fmt.Sprint(t.ID),
			Date:    t.Date,
			Title:   t.Title,
			Comment: t.Comment,
			Repeat:  t.Repeat,
		}
	}

	writeJSON(w, map[string]any{"tasks": respTasks}) // ✅ REPLACED ORIGINAL WRITEJSON
}
