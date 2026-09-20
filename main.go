package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"todo_app/pkg/api"
	"todo_app/pkg/db"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db" // дефолт
	}

	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка при инициализации БД: %v", err)
	}
	defer db.Close()

	api.Init()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	webDir := "./web"

	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	addr := "0.0.0.0:" + port

	fmt.Printf("Сервер запущен на http://localhost:%s\n", port)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске: %v", err)
	}
}
