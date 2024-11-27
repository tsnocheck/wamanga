package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("34r*#*FDEC") // Секретный ключ для подписи

type Claims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}

// Функция для генерации JWT токена
func GenerateJWT(id uint) (string, error) {
	// Время истечения срока действия токена
	expirationTime := time.Now().Add(30 * 24 * time.Hour)

	// Создание объекта claims
	claims := &Claims{
		Id: fmt.Sprintf("%d", id), // Преобразуем id в строку
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // Время истечения
		},
	}

	// Генерация токена с HMAC подписью
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

func DecodeToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	log.Println("JWT Key:", jwtKey)

	// Парсинг токена
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})

	// Обработка ошибок при парсинге токена
	if err != nil {
		log.Println("Error parsing token:", err)
		return nil, err
	}

	// Проверка на валидность токена
	if !token.Valid {
		log.Println("Invalid token")
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
