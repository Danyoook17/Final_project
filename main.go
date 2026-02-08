package main

import (
	"log"
	"net/http"
	"os"
	"todoapp/pkg/api"
	"todoapp/pkg/db"
)

func main() {
	// Порт
	port := "7540"
	if p := os.Getenv("TODO_PORT"); p != "" {
		port = p
	}

	// Путь к БД
	dbFile := "scheduler.db"
	if v := os.Getenv("TODO_DBFILE"); v != "" {
		dbFile = v
	}

	// Инициализация БД
	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}
	api.Init()

	// Раздача фронта
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Println("Server started on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

