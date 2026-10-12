package main

import (
	"final/pkg/db"
	"final/pkg/server"
	"fmt"
	"os"
)

const (
	defaultPort   = "7540"
	defaultDBFile = "scheduler.db"
)

func main() {
	// 1. Определяем путь к БД
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// 2. Инициализируем БД
	err := db.Init(dbFile)
	if err != nil {
		fmt.Println("Ошибка БД:", err)
		return
	}
	defer db.DB.Close()

	// 3. Запуск сервера
	if err := server.Run(defaultPort); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}

}
