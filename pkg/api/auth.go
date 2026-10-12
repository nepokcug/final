package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("sign")
var todoPass string

func initAuth() {
	todoPass = os.Getenv("TODO_PASSWORD")
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, map[string]string{"error": "метод не поддерживается"}, http.StatusMethodNotAllowed)
	}
	// Читаем пароль из тела
	var creds struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeJson(w, map[string]string{"error": "неверный JSON"}, http.StatusBadRequest)
		return
	}

	// Сверяем пароль
	if creds.Password != todoPass {
		writeJson(w, map[string]string{"error": "неверный пароль"}, http.StatusUnauthorized)
		return
	}

	// Хэш пароля
	hash := sha256.Sum256([]byte(todoPass))
	hashStr := hex.EncodeToString(hash[:])

	// Создание JWT токена
	claims := jwt.MapClaims{
		"hash": hashStr,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка токена"}, http.StatusInternalServerError)
		return
	}

	// Возвращаем токен
	writeJson(w, map[string]string{"token": signed}, http.StatusOK)
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Есть ли пароль
		if len(todoPass) == 0 {
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
		if !validateJWT(cookie.Value, todoPass) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		// Все ок -> вызываем обработчик
		next(w, r)
	}
}
