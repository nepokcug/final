package server

import (
	"fmt"
	"net/http"
	"os"

	"final/pkg/api"
)

func Run(defaultPort string) error {

	api.Init()

	// 3. Запускаем сервер
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))
	fmt.Println("Сервер запущен на http://localhost:" + port)
	return http.ListenAndServe(":"+port, nil)
}
