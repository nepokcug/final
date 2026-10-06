package main

import (
	"final/pkg/server"
	"fmt"
)

func main() {

	if err := server.Run(); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
}
