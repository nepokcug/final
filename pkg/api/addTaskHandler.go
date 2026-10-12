package api

import (
	"encoding/json"
	"final/pkg/db"
	"net/http"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, map[string]string{"error": "метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}
	// 1. Десериализация JSON → Task
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "неверный JSON"}, http.StatusBadRequest)
		return
	}

	// 2. Проверка title
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "отсутвует заголовок"}, http.StatusBadRequest)
		return
	}

	// 3. Проверка даты
	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// 4. Добавление в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// 5. Ответ с ID
	writeJson(w, map[string]int64{"id": id}, http.StatusCreated)
}
