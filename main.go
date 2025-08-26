package main

import (
	"log"
	"os"

	"LAST_TODO_2/pkg/api" //добавил импорт для InitAuth
	"LAST_TODO_2/pkg/db"
	"LAST_TODO_2/pkg/server"
)

func main() {
	// читаем порт из переменной окружения
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// читаем путь к базе данных из переменной окружения
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// инициализация базы данных
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	defer db.Close() // закрываем базу при выходе

	api.InitAuth() //инициализация авторизации на старте

	// создаём и запускаем сервер
	s := server.New(port, "./web")
	log.Printf("Server starting at http://localhost:%s ...", port)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
