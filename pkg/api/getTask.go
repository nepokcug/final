package api

import (
	"final/pkg/db"
	"net/http"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, map[string]string{"error": "метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}
	// 1. Получаем id
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "не указан id"}, http.StatusOK)
		return
	}
	// 2. Получаем задачу
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}

	// 3. Отвечаем
	writeJson(w, task, http.StatusOK)
}
