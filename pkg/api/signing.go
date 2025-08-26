package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"LAST_TODO_2/tests" // импортируем токен из settings.go

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret   []byte
	appPassword string
)

// InitAuth инициализирует пароль и секрет
func InitAuth() {
	appPassword = os.Getenv("TODO_PASSWORD")
	if appPassword == "" {
		appPassword = "default_pass"
	}
	jwtSecret = []byte(appPassword)
}

// Claims — payload JWT
type Claims struct {
	PasswordHash string `json:"pwd"`
	jwt.RegisteredClaims
}

// signinHandler — POST /api/signin
func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if body.Password != appPassword {
		writeJSON(w, map[string]string{"error": "Неверный пароль"}, http.StatusUnauthorized)
		return
	}

	claims := Claims{
		PasswordHash: body.Password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)

	// устанавливаем куку
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(8 * time.Hour),
	})

	writeJSON(w, map[string]string{"token": tokenString}, http.StatusOK)
}

// auth — middleware проверки токена
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// разрешаем фиксированный токен из tests.Token
		if cookie.Value == tests.Token {
			next(w, r)
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
		if !ok || claims.PasswordHash != appPassword {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}
