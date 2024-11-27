package controllers

import (
	"github.com/gin-gonic/gin"
	"main/src/models"
	"net/http"
	"strings"
)

// Login godoc
// @Summary Логин пользователя
// @Description Логин с использованием email и пароля
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Данные для входа"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error(), "data": nil})
		return
	}

	authResponse, err := input.Login()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Login successful", "data": authResponse})
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Регистрация нового пользователя с использованием JSON данных
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Данные для регистрации пользователя"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error(), "data": nil})
		return
	}

	// Вызов метода Register из модели для регистрации пользователя
	authResponse, err := input.Register()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error(), "data": nil})
		return
	}

	// Возвращаем успешный ответ с токеном и данными пользователя
	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "User registered successfully", "data": authResponse})
}

// GetProfile godoc
// @Summary Получить профиль пользователя
// @Description Получение данных профиля пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security apiKey  // Указывает, что требуется API ключ
// @Param Authorization header string true "API Key in Bearer format"  // Указываем, что ключ передается в заголовке в формате Bearer
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/profile [get]
func GetProfile(c *gin.Context) {
	// Извлекаем токен из заголовка Authorization
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Authorization header is missing", "data": nil})
		return
	}

	// Убедимся, что токен начинается с "Bearer "
	if !strings.HasPrefix(token, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Invalid token format. Bearer token required", "data": nil})
		return
	}

	// Извлекаем сам токен (удаляем префикс "Bearer ")
	reqToken := token[len("Bearer "):]

	// Декодируем и проверяем токен
	claims, err := models.DecodeToken(reqToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Invalid token: " + err.Error(), "data": nil})
		return
	}

	// Получаем пользователя на основе ID из токена
	var user models.User
	err = models.Database.First(&user, claims.Id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "User not found", "data": nil})
		return
	}

	// Формируем ответ с данными пользователя
	authResponse := models.AuthResponse{
		User:  &user,
		Token: reqToken, // При необходимости добавляем токен в ответ
	}

	// Возвращаем успешный ответ с данными пользователя
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User profile fetched successfully", "data": authResponse})
}
