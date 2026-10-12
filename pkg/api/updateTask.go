package api

import (
	"encoding/json"
	"final/pkg/db"
	"net/http"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJson(w, map[string]string{"error": "метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}
	// 1. Десериализация
	task := &db.Task{}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "неверный JSON"}, http.StatusBadRequest)
		return
	}
	// 2. Проверка id
	if task.ID == "" {
		writeJson(w, map[string]string{"error": "не указан id"}, http.StatusBadRequest)
		return
	}
	// 3. Проверка Title
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "не указан заголовок"}, http.StatusBadRequest)
		return
	}
	// 4. Проверка даты
	if err := checkDate(task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	// 5. Обновление
	if err := db.UpdateTask(task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	// 6. Ответ
	writeJson(w, map[string]string{}, http.StatusOK)
}
