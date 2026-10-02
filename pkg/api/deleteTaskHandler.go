package api

import (
	"final/pkg/db"
	"net/http"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Ищем id
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "не указан id"})
		return
	}
	// 1. Удаляем
	if err := db.DeleteTask(id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}
