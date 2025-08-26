// pkg/api/api.go
package api

import (
	"net/http"
	"time"
)

// Init подключает все обработчики API
func Init() {
	// публичные эндпоинты
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/nextdate", nextDateHandler)

	// защищённые эндпоинты через auth middleware
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/task/done", auth(doneHandler))

}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(isoDate, nowStr)
		if err != nil {
			http.Error(w, ErrInvalidDate.Error(), http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(result))
}
