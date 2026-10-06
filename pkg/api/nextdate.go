package api

import (
	"fmt"
	"net/http"
	"time"

	"final/pkg/nextdate"
)

// Формат даты
const dateFormat = "20060102"

// nextDayHandler обрабатывает GET /api/nextdate
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Читаем параметры
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// 2. Если now не указан — берём текущую дату
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "неверный формат now", http.StatusBadRequest)
			return
		}
	}

	// 3. Вызываем NextDate
	next, err := nextdate.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 4. Возвращаем результат
	fmt.Fprint(w, next)
}
