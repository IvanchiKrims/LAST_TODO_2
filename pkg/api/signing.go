package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Секрет для подписи токена (можно использовать TODO_PASSWORD)
var jwtSecret = []byte("secret_for_signing") // можно динамически поставить из TODO_PASSWORD

// payload токена
type Claims struct {
	PasswordHash string `json:"pwd"`
	jwt.RegisteredClaims
}

// signinHandler — POST /api/signin
func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]any{"error": "method not allowed"})
		return
	}

	var body struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, map[string]any{"error": "invalid JSON"})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		writeJSON(w, map[string]any{"error": "password not set"})
		return
	}

	if body.Password != pass {
		writeJSON(w, map[string]any{"error": "Неверный пароль"})
		return
	}

	// создаём токен JWT
	claims := Claims{
		PasswordHash: body.Password, // можно заменить на хэш
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		writeJSON(w, map[string]any{"error": "failed to generate token"})
		return
	}

	// устанавливаем куку на 8 часов
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(8 * time.Hour),
	})

	writeJSON(w, map[string]any{"token": tokenString})
}

// auth — middleware для проверки JWT из куки
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(t *jwt.Token) (any, error) {
				return jwtSecret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*Claims)
			if !ok || claims.PasswordHash != pass {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	})
}
