package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"main/src/models/structur"
	"net/http"
	"strconv"
)

// GetComicsInfo godoc
// @Summary Получить информацию о комиксе
// @Description Получить детальную информацию о комиксе по имени
// @Tags Comics
// @Accept json
// @Produce json
// @Security Name  // Указывает, что требуется API ключ
// @Param name query string true "Название комикса"  // Имя комикса передается через query параметр
// @Success 200 {object} structur.Comics
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /comics/info [get]
func GetComicsInfo(c *gin.Context) {
	// Извлекаем имя комикса через query параметр
	comicName := c.DefaultQuery("name", "")
	if comicName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Name parameter is missing", "data": nil})
		return
	}

	// Получаем информацию о комиксе по имени
	comicInfo, err := structur.GetComicsInfo(comicName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Comic not found", "data": nil})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error(), "data": nil})
		}
		return
	}

	// Формируем успешный ответ с информацией о комиксе
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Get info successful", "data": comicInfo})
}

// DeleteComics godoc
// @Summary Удалить комикс
// @Description Удаление комикса
// @Tags Comics
// @Accept json
// @Produce json
// @Security Name  // Указывает, что требуется API ключ
// @Param name query string true "Название комикса"  // Имя комикса передается через query параметр
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /comics/delete [delete]
func DeleteComics(c *gin.Context) {
	// Извлекаем имя комикса через query параметр
	comicName := c.DefaultQuery("name", "")
	if comicName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Name parameter is missing", "data": nil})
		return
	}

	// Проверяем, существует ли комикс в базе данных
	comicInfo, err := structur.GetComicsInfo(comicName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Comic not found", "data": nil})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error(), "data": nil})
		}
		return
	}

	// Удаляем комикс из базы данных
	comicInfo, err = structur.DeleteComicByName(comicName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to delete comic", "data": nil})
		return
	}

	// Формируем успешный ответ после удаления
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Comic deleted successfully", "data": comicInfo})
}

// CreateComics godoc
// @Summary Создать новый комикс
// @Description Создание нового комикса с детальной информацией, включая изображение обложки и баннера.
// @Tags Comics
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Название комикса"
// @Param alternative_name formData string true "Альтернативное название комикса"
// @Param description formData string true "Описание комикса"
// @Param rating formData float32 true "Рейтинг комикса"
// @Param image_path formData file true "Изображение обложки комикса"
// @Param banner_path formData file true "Изображение баннера комикса"
// @Param type_comics formData string true "Тип комикса"
// @Param author formData string true "Автор комикса"
// @Param original_author formData string true "Оригинальный автор комикса"
// @Param year formData int true "Год выпуска комикса"
// @Param is_finished formData bool true "Статус завершенности комикса"
// @Param pegi formData string true "Возрастной рейтинг"
// @Param status formData string true "Статус комикса"
// @Param transfer_status formData string true "Статус переноса комикса"
// @Param views formData int32 true "Количество просмотров"
// @Param likes formData int32 true "Количество лайков"
// @Param hidden formData bool true "Статус скрытости комикса"
// @Param tags formData array true "Теги комикса"  items({"type": "string"}) // Ожидается массив строк
// @Param genres formData array true "Жанры комикса"  items({"type": "string"}) // Ожидается массив строк
// @Param bookmark formData int true "Закладки"
// @Success 200 {object} structur.Comics
// @Failure 400 {object} map[string]interface{}
// @Router /comics/create [post]
func CreateComics(c *gin.Context) {
	var comic structur.Comics

	// Присваивание полей вручную с учетом типов
	comic.Name = c.PostForm("name")
	comic.AlternativeName = c.PostForm("alternative_name")
	comic.Description = c.PostForm("description")

	// Преобразование и присвоение числовых полей
	if rating, err := strconv.ParseFloat(c.PostForm("rating"), 32); err == nil {
		comic.Rating = float32(rating)
	}
	if year, err := strconv.Atoi(c.PostForm("year")); err == nil {
		comic.Year = year
	}
	if views, err := strconv.Atoi(c.PostForm("views")); err == nil {
		comic.Views = int32(views)
	}
	if likes, err := strconv.Atoi(c.PostForm("likes")); err == nil {
		comic.Likes = int32(likes)
	}
	if bookmark, err := strconv.Atoi(c.PostForm("bookmark")); err == nil {
		comic.Bookmark = bookmark
	}

	// Присваивание остальных полей
	comic.Type = structur.ComicsType(c.PostForm("type_comics"))
	comic.Author = c.PostForm("author")
	comic.Artist = c.PostForm("original_author")
	comic.IsFinished = c.PostForm("is_finished") == "true"
	comic.Pegi = structur.PegiType(c.PostForm("pegi"))
	comic.Status = structur.StatusType(c.PostForm("status"))
	comic.TransferStatus = structur.StatusType(c.PostForm("transfer_status"))
	comic.Hidden = c.PostForm("hidden") == "true"
	comic.Tags = pq.StringArray(c.PostFormArray("tags"))
	comic.Genres = pq.StringArray(c.PostFormArray("genres"))

	// Получение файлов изображения
	coverFile, err := c.FormFile("image_path")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Cover image is required", "data": nil})
		return
	}
	bannerFile, err := c.FormFile("banner_path")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Banner image is required", "data": nil})
		return
	}

	// Открываем файлы изображения для чтения
	coverStream, err := coverFile.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to open cover image", "data": nil})
		return
	}
	defer coverStream.Close()

	bannerStream, err := bannerFile.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Failed to open banner image", "data": nil})
		return
	}
	defer bannerStream.Close()

	// Вызов функции UploadComics для сохранения комикса и изображений
	comicResponse, err := structur.UploadComics(&comic, coverStream, bannerStream)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error(), "data": nil})
		return
	}

	// Ответ с успешным созданием комикса
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Comic created successfully", "data": comicResponse})
}
