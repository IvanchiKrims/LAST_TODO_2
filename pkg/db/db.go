package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var (
	// Глобальная переменная для хранения открытой базы
	DB *sql.DB
)

// SQL-схема для создания таблицы и индекса
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT '',
	title VARCHAR(128) NOT NULL DEFAULT '',
	comment TEXT NOT NULL DEFAULT '',
	repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`

// Init открывает базу данных и создаёт таблицу, если её нет
func Init(dbFile string) error {
	// проверка переменной окружения TODO_DBFILE
	if envFile := os.Getenv("TODO_DBFILE"); envFile != "" {
		dbFile = envFile
	}

	// проверяем существование файла
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		install = true
	}

	// открываем базу данных
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть базу данных: %w", err)
	}

	// если файла не было, создаём таблицу и индекс
	if install {
		if _, err := DB.Exec(schema); err != nil {
			return fmt.Errorf("не удалось создать таблицу: %w", err)
		}
	}

	return nil
}

// Close закрывает соединение с базой данных
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
