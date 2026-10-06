package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("sign")

func signInHandler(w http.ResponseWriter, r *http.Request) {
	// Читаем пароль из тела
	var creds struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeJson(w, map[string]string{"error": "неверный JSON"})
		return
	}

	// Сверяем пароль
	pass := os.Getenv("TODO_PASSWORD")
	if creds.Password != pass {
		writeJson(w, map[string]string{"error": "неверный пароль"})
		return
	}

	// Хэш пароля
	hash := sha256.Sum256([]byte(pass))
	hashStr := hex.EncodeToString(hash[:])

	// Создание JWT токена
	claims := jwt.MapClaims{
		"hash": hashStr,
		"exp":  time.Now().Add(time.Hour * 8).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка токена"})
		return
	}

	// Возвращаем токен
	writeJson(w, map[string]string{"token": signed})
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Есть ли пароль
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) == 0 {
			next(w, r)
			return
		}

		// Читаем куку "token"
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		// Проверяем токен
		if !validateJWT(cookie.Value, pass) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		// Все ок -> вызываем обработчик
		next(w, r)
	}
}
