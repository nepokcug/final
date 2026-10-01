package api

import (
	"final/pkg/db"
	"net/http"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем id
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "не указан id"})
		return
	}
	// 2. Получаем задачу
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	// 3. Отвечаем
	writeJson(w, task)
}
