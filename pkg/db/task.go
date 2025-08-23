package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

// Task — модель задачи в БД
type Task struct {
	ID      int64  `db:"id" json:"id"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

// AddTask — добавляет новую задачу
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetTask — возвращает задачу по ID
func GetTask(id string) (*Task, error) {
	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}

	var t Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err = DB.QueryRow(query, taskID).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return &t, nil
}

// UpdateTask — обновляет все поля задачи
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}
	return nil
}

// UpdateDate — обновляет только дату задачи
func UpdateDate(next string, id string) error {
	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	query := `UPDATE scheduler SET date=? WHERE id=?`
	res, err := DB.Exec(query, next, taskID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// DeleteTask — удаляет задачу по ID
func DeleteTask(id string) error {
	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid id")
	}
	query := `DELETE FROM scheduler WHERE id=?`
	res, err := DB.Exec(query, taskID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

// Tasks — возвращает ближайшие n задач
func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat 
	          FROM scheduler 
	          ORDER BY date 
	          LIMIT ?`
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}

// TasksByDate — возвращает задачи по конкретной дате
func TasksByDate(date string, limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat 
	          FROM scheduler 
	          WHERE date=? 
	          ORDER BY id 
	          LIMIT ?`
	rows, err := DB.Query(query, date, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}

// TasksBySearch — ищет задачи по подстроке в title или comment
func TasksBySearch(search string, limit int) ([]*Task, error) {
	like := "%" + search + "%"
	query := `SELECT id, date, title, comment, repeat 
	          FROM scheduler 
	          WHERE title LIKE ? OR comment LIKE ? 
	          ORDER BY date 
	          LIMIT ?`
	rows, err := DB.Query(query, like, like, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, nil
}
