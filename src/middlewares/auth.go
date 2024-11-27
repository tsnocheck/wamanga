package middlewares

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"main/src/models"
)

// AuthMiddleware проверяет наличие и корректность API ключа в заголовке Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization
		var token = c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is missing"})
			return
		}

		// Проверяем, что токен начинается с "Bearer "
		const bearerPrefix = "Bearer "
		splitToken := strings.Split(token, bearerPrefix)

		// Проверяем, что строка разделена корректно
		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token format. Bearer token required"})
			return
		}

		// Извлекаем сам токен
		reqToken := splitToken[1]

		// Декодируем и проверяем токен
		claims, err := models.DecodeToken(reqToken)
		if err != nil {
			log.Printf("Error decoding token: %v\n", err) // Логирование ошибки при декодировании токена
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		// Логируем данные claims для отладки
		log.Printf("Claims: %v\n", claims)

		// Используем данные из claims для установки userId в контексте
		c.Set("userId", claims.Id) // Замените на соответствующее поле, если нужно

		// Переходим к следующему обработчику
		c.Next()
	}
}
