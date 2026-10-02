package api

import (
	"final/pkg/db"
	"net/http"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Ищем id
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "отсутствует id"})
		return
	}
	// 2. Выполняем (там же если без правила удалится)
	if err := db.DoneTask(id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]string{})
}
