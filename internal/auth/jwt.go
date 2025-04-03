package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey []byte
	adminKey  []byte
)

type Claims struct {
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	PCID    string `json:"pc_id,omitempty"` // Только для пользовательских токенов
	jwt.RegisteredClaims
}

// Инициализация ключей
func Init(secret, adminSecret string) {
	secretKey = []byte(secret)
	adminKey = []byte(adminSecret)
}

// Генерация токена
func GenerateToken(userID string, isAdmin bool, pcID string) (string, error) {
	claims := &Claims{
		UserID:  userID,
		IsAdmin: isAdmin,
		PCID:    pcID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	key := secretKey
	if isAdmin {
		key = adminKey
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

// Валидация токена
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		if claims.IsAdmin {
			return adminKey, nil
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
