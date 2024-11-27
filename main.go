package main

import (
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "main/docs"
	"main/src/models"
	"main/src/models/structur"
	"main/src/routes"
	"main/src/utils"
)

func main() {
	// Загрузка переменных окружения
	utils.LoadEnv()

	// Подключение к базе данных и миграции
	models.OpenDatabaseConnection()
	models.AutoMigrateModels()
	structur.AutoMigrateComics()

	r := routes.SetupRoutes()

	// Настройка Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера
	r.Run(":8080") // Слушаем порт 8080
}
