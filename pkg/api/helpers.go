package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"final/pkg/db"
	"final/pkg/nextdate"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func checkDate(task *db.Task) error {
	// Если дата не указана присваиваем текущее время
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}
	// Проверяем что указана корректная дата
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return errors.New("неверный формат даты")
	}
	// Если есть правило вычисляем следующую дату
	var next string
	if task.Repeat != "" {
		next, err = nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if nextdate.AfterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			task.Date = next
		}
	}
	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func validateJWT(tokenStr, pass string) bool {
	// Парсим токен
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return false
	}

	// Достаем хэш
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	hashFromToken, ok := claims["hash"].(string)
	if !ok {
		return false
	}

	// Хэш текущего пароля
	hashNow := sha256.Sum256([]byte(pass))
	hashNowStr := hex.EncodeToString(hashNow[:])

	// Сравниваем
	return hashNowStr == hashFromToken
}
