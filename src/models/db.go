package models

import (
	"fmt"
	"main/src/utils"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func OpenDatabaseConnection() {
	var err error

	host := os.Getenv("POSTGRES_HOST")
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	databaseName := os.Getenv("POSTGRES_DATABASE")
	port := os.Getenv("POSTGRES_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, username, password, databaseName, port)

	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		// Получаем логгер из utils и записываем ошибку
		utils.Logger("Error connecting to database", "error", err)
		panic(err)
	} else {
		utils.Logger("Connected to database", "info")
		fmt.Println("🚀🚀🚀---ASCENDE SUPERIUS---🚀🚀🚀")
	}
}

func AutoMigrateModels() {
	Database.AutoMigrate(&User{})
}
