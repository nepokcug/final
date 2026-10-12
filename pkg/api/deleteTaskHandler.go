package api

import (
	"final/pkg/db"
	"net/http"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJson(w, map[string]string{"error": "метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}
	// 1. Ищем id
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "не указан id"}, http.StatusNotFound)
		return
	}
	// 1. Удаляем
	if err := db.DeleteTask(id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}

	writeJson(w, map[string]string{}, http.StatusOK)
}
