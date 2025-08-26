package server

import (
	"fmt"
	"net/http"

	"LAST_TODO_2/pkg/api"
)

type Server struct {
	port   string
	webDir string
}

func New(port string, webDir string) *Server {
	return &Server{port: port, webDir: webDir}
}

func (s *Server) Start() error {
	// статический файловый сервер (фронтенд)
	fs := http.FileServer(http.Dir(s.webDir))
	http.Handle("/", fs)

	// инициализируем API
	api.Init()

	// запуск сервера
	fmt.Printf("Server starting at http://localhost:%s ...\n", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), nil)
}
