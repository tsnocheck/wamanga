package routes

import (
	"main/src/controllers"
	"main/src/middlewares"

	"github.com/gin-gonic/gin"
)

// startupsGroupRouter - настройка маршрутов для авторизации
func startupsGroupRouter(baseRouter *gin.RouterGroup) {
	auth := baseRouter.Group("/auth")

	// Роуты авторизации
	auth.GET("/profile", middlewares.AuthMiddleware(), controllers.GetProfile)
	auth.POST("/login", controllers.Login)
	auth.POST("/register", controllers.Register) // <---- этот маршрут должен существовать
}

func zalupaCom(baseRouter *gin.RouterGroup) {
	auth := baseRouter.Group("/comics")

	auth.GET("/info", controllers.GetComicsInfo)
	auth.POST("/create", controllers.CreateComics)
}

// SetupRoutes - настройка всех маршрутов
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Группируем версии API
	apiV1 := r.Group("/api/v1")

	// Добавляем маршруты для авторизации
	startupsGroupRouter(apiV1)
	zalupaCom(apiV1)

	return r
}
