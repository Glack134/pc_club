package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Секретный ключ (должен совпадать с серверным)
	secret := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDM0NTE2NDMsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjoiYWRtaW4ifQ.zj_ZukYHGf7v58LtOGx7qlGRhEpepAxrKXc6s4saPFc"

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
