package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"final/pkg/api"
	"final/pkg/db"
	"final/tests"
)

func main() {
	// 1. Определяем путь к БД
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = tests.DBFile
	}

	// 2. Инициализируем БД
	err := db.Init(dbFile)
	if err != nil {
		fmt.Println("Ошибка БД:", err)
		return
	}

	api.Init()

	// 3. Запускаем сервер
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(tests.Port)
	}

	http.Handle("/", http.FileServer(http.Dir("../web")))
	fmt.Println("Сервер запущен на http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
