package structur

import (
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"io"
	"log"
	"main/src/models"
	"os"
	"path/filepath"
	"time"
)

type ComicsType string
type PegiType string
type StatusType string

// swagger:ignore
type StringArray pq.StringArray

const (
	Manga      ComicsType = "Manga"
	Manhva     ComicsType = "Webtoon"
	Manhua     ComicsType = "Webtoon"
	Comic      ComicsType = "Comic"
	LifeComic  ComicsType = "Life Comic"
	WebComic   ComicsType = "Web Comic"
	Manuscript ComicsType = "Manuscript"
)

const (
	Pegi3  PegiType = "3+"
	Pegi6  PegiType = "6+"
	Pegi12 PegiType = "12+"
	Pegi16 PegiType = "16+"
	Pegi18 PegiType = "18+"
)

const (
	Started    StatusType = "В процессе"
	isFinished StatusType = "Окончено"
	Paused     StatusType = "Приостановлено"
	Abandoned  StatusType = "Заброшено"
	Аnnounced  StatusType = "Анонсировано"
)

// TODO: я хз похуй мне на это
type Comics struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" `
	AlternativeName string         `json:"alternative_name" ` // Необязательное поле
	Description     string         `json:"description" `
	Rating          float32        `json:"rating" `
	ImagePath       string         `json:"image_path" `
	BannerPath      string         `json:"banner_path" `
	Type            ComicsType     `json:"type_comics" `
	Author          string         `json:"author" `
	Artist          string         `json:"original_author" `
	Year            int            `json:"year" `
	IsFinished      bool           `json:"is_finished" `
	Pegi            PegiType       `json:"pegi" `
	Status          StatusType     `json:"status" `
	TransferStatus  StatusType     `json:"transfer_status" `
	Views           int32          `json:"views" `
	Likes           int32          `json:"likes" `
	Hidden          bool           `json:"hidden" `
	PublishedOn     time.Time      `json:"published_on"` // Необязательное поле
	UpdatedAt       time.Time      `json:"updated_at"`   // Необязательное поле
	Tags            pq.StringArray `json:"tags" gorm:"type:text[]" swaggertype:"array,string" `
	Genres          pq.StringArray `json:"genres" gorm:"type:text[]" swaggertype:"array,string" `
	Bookmark        int            `json:"bookmark"`
}

func GetComicsInfo(name string) (*Comics, error) {
	var comics Comics
	log.Printf("Received name: %s", name) // Добавляем вывод значения name

	err := models.Database.Where("name = ?", name).First(&comics).Error
	if err != nil {
		log.Printf("Error retrieving comic: %v", err)
		return &Comics{}, err
	}

	log.Printf("Found comic: %+v", comics)
	return &comics, nil
}

func UploadComics(dto *Comics, imageStream, bannerStream io.Reader) (*Comics, error) {
	var existingComic Comics

	if err := models.Database.Where("alternative_name = ?", dto.AlternativeName).First(&existingComic).Error; err == nil {
		// Если запись существует, возвращаем ошибку
		return nil, errors.New("a comic with the same alternative name already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Если возникла другая ошибка при поиске, возвращаем её
		return nil, fmt.Errorf("failed to check existing comic: %w", err)
	}

	// Логирование для отладки
	log.Printf("Checking for existing comic with name: %s", dto.Name)

	// Создание AlternativeName на английском, если он пуст или содержит кириллицу
	if dto.AlternativeName == "" || containsCyrillic(dto.AlternativeName) {
		translatedName := slug.Make(dto.Name)
		dto.AlternativeName = translatedName
	}

	// Установка текущих даты и времени для PublishedOn и UpdatedAt, если они не заданы
	now := time.Now()
	if dto.PublishedOn.IsZero() {
		dto.PublishedOn = now
	}
	if dto.UpdatedAt.IsZero() {
		dto.UpdatedAt = now
	}

	// Создание директории для изображений, если она еще не существует
	imageDir := fmt.Sprintf("./main/images/%s", dto.AlternativeName)
	if err := os.MkdirAll(filepath.Join(imageDir, "banners"), os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create banners directory: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(imageDir, "cover"), os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create cover directory: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(imageDir, "chapters"), os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create chapters directory: %w", err)
	}

	var err error

	dto.ImagePath, err = saveImage(imageStream, imageDir, "cover.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to save cover image: %w", err)
	}

	dto.BannerPath, err = saveImage(bannerStream, imageDir, "banner.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to save banner image: %w", err)
	}

	// Сохранение в базу данных
	err = models.Database.Create(dto).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

// saveImage сохраняет изображение из потока в файл и возвращает путь
func saveImage(imageStream io.Reader, dir, filename string) (string, error) {
	path := filepath.Join(dir, filename)
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, imageStream)
	if err != nil {
		return "", err
	}
	return path, nil
}

// Проверка, содержит ли строка кириллические символы
func containsCyrillic(text string) bool {
	for _, r := range text {
		if r >= 'А' && r <= 'я' {
			return true
		}
	}
	return false
}

func DeleteComicByName(name string) (*Comics, error) {
	var comic Comics

	// Поиск комикса по альтернативному имени
	err := models.Database.Where("alternative_name = ?", name).First(&comic).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("comic not found")
		}
		return nil, err
	}

	// Удаление найденного комикса
	if err := models.Database.Delete(&comic).Error; err != nil {
		return nil, fmt.Errorf("failed to delete comic: %w", err)
	}

	// Удаление папки с изображениями
	imageDir := fmt.Sprintf("./main/images/%s", comic.AlternativeName)
	if err := os.RemoveAll(imageDir); err != nil {
		return nil, fmt.Errorf("failed to delete image directory: %w", err)
	}

	// Возврат удалённого комикса
	return &comic, nil
}

func UpdateComicsInfo(name string, updatedFields map[string]interface{}, newCover, newBanner io.Reader) (*Comics, error) {
	var comics Comics
	err := models.Database.Where("alternative_name = ?", name).First(&comics).Error
	if err != nil {
		return nil, err
	}

	// Update cover image if newCover is provided
	if newCover != nil {
		imageDir := fmt.Sprintf("./main/images/%s", comics.AlternativeName)
		comics.ImagePath, err = saveImage(newCover, imageDir, "cover.jpg")
		if err != nil {
			return nil, fmt.Errorf("failed to save new cover image: %w", err)
		}
		updatedFields["image_path"] = comics.ImagePath
	}

	// Update banner image if newBanner is provided
	if newBanner != nil {
		imageDir := fmt.Sprintf("./main/images/%s", comics.AlternativeName)
		comics.BannerPath, err = saveImage(newBanner, imageDir, "banner.jpg")
		if err != nil {
			return nil, fmt.Errorf("failed to save new banner image: %w", err)
		}
		updatedFields["banner_path"] = comics.BannerPath
	}

	err = models.Database.Model(&comics).Updates(updatedFields).Error
	if err != nil {
		return nil, err
	}
	return &comics, nil
}

func AutoMigrateComics() {
	models.Database.AutoMigrate(&Comics{})
}
