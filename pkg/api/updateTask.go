package api

import (
	"encoding/json"
	"final/pkg/db"
	"net/http"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Десериализация
	task := &db.Task{}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "неверный JSON"})
		return
	}
	// 2. Проверка id
	if task.ID == "" {
		writeJson(w, map[string]string{"error": "не указан id"})
		return
	}
	// 3. Проверка Title
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "не указан заголовок"})
		return
	}
	// 4. Проверка даты
	if err := checkDate(task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	// 5. Обновление
	if err := db.UpdateTask(task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	// 6. Ответ
	writeJson(w, map[string]string{})
}
