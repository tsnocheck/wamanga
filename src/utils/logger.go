package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

func Logger(errorMessage string, typeError string, err ...error) {
	log := logrus.New()

	// Открытие файла для записи логов
	file, fileErr := os.OpenFile("./src/logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if fileErr != nil {
		log.Fatal("Failed to open log file", fileErr)
	}

	log.SetOutput(file)

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Настройка уровня логирования
	log.SetLevel(logrus.DebugLevel)

	// Логирование в зависимости от типа ошибки
	switch typeError {
	case "info":
		if len(err) > 0 {
			log.Info(errorMessage, err[0])
		} else {
			log.Info(errorMessage)
		}
	case "warn":
		if len(err) > 0 {
			log.Warn(errorMessage, err[0])
		} else {
			log.Warn(errorMessage)
		}
	case "error":
		if len(err) > 0 {
			log.Error(errorMessage, err[0])
		} else {
			log.Error(errorMessage)
		}
	default:
		if len(err) > 0 {
			log.Debug("Unknown log level: ", errorMessage, err[0])
		} else {
			log.Debug("Unknown log level: ", errorMessage)
		}
	}
}
