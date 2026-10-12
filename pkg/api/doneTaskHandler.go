package api

import (
	"final/pkg/db"
	"final/pkg/nextdate"
	"net/http"
	"time"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, map[string]string{"error": "метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}
	// 1. Ищем id
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "отсутствует id"}, http.StatusBadRequest)
		return
	}
	// 2. Получаем задачу
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}
	// 3. Если правило пустое -> удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJson(w, map[string]string{}, http.StatusOK)
		return
	}
	// 4. Иначе → вычисляем next
	next, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// 5. Обновляем дату
	if err := db.UpdateDate(next, id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{}, http.StatusOK)
}
