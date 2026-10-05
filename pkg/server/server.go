package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"final/pkg/api"
	"final/pkg/db"
	"final/tests"
)

func Run() error {
	// 1. Определяем путь к БД
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// 2. Инициализируем БД
	err := db.Init(dbFile)
	if err != nil {
		return fmt.Errorf("ошибка БД: %w", err)
	}

	api.Init()

	// 3. Запускаем сервер
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))
	fmt.Println("Сервер запущен на http://localhost:" + port)
	return http.ListenAndServe(":"+port, nil)
}
