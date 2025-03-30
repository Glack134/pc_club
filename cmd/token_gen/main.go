package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Секретный ключ (должен совпадать с серверным)
	secret := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDM0NDEzODAsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjoiYWRtaW4ifQ.yMFyqjXqCUWt0oA9nJ7o0tZbq5t3YUaMOo3Mx6qydZU" // Замените на реальный секрет

	// Создаем claims (данные токена)
	claims := jwt.MapClaims{
		"user_id":  "admin",
		"is_admin": true,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	// Генерируем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal("Failed to generate token:", err)
	}

	fmt.Println("Generated token:")
	fmt.Println(tokenString)
}
