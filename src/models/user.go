package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"regexp"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	ID     uint   `json:"id" gorm:"primaryKey"`          // Добавляем ID для AuthResponse
	UserID uint   `json:"user_id"`                       // Внешний ключ для пользователя
	User   *User  `json:"user" gorm:"foreignKey:UserID"` // Связь с моделью User через внешний ключ
	Token  string `json:"token"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Fetches a user from the database
func FetchUser(id uint) (*User, error) {
	var user User
	err := Database.Where("id = ?", id).First(&user).Error
	if err != nil {
		return &User{}, err
	}
	return &user, nil
}

func FetchUserByEmail(email string) User {
	var userFromDb User
	Database.Where("email = ?", email).First(&userFromDb)

	return userFromDb
}

func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(hashedPassword)
	return err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (user *User) Register() (*AuthResponse, error) {
	// Проверяем, что email имеет корректный формат
	if !emailRegex.MatchString(user.Email) {
		return nil, errors.New("invalid email format")
	}

	// Проверяем, существует ли уже пользователь с таким email
	var userFromDb User
	err := Database.Where("email = ?", user.Email).First(&userFromDb).Error

	// Если пользователь найден, возвращаем ошибку
	if err == nil {
		return nil, errors.New("email already taken")
	}

	err = Database.Where("username = ?", user.Username).First(&userFromDb).Error

	if err == nil {
		return nil, errors.New("username already taken")
	}

	// Если ошибка не связана с тем, что пользователь не найден (например, другая ошибка)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Хешируем пароль
	err = user.HashPassword()
	if err != nil {
		return nil, err
	}

	// Сохраняем пользователя в базе данных
	err = Database.Create(user).Error
	if err != nil {
		return nil, err
	}

	// Генерируем JWT токен
	token, err := GenerateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	// Формируем ответ
	response := AuthResponse{
		User:  user, // Обратите внимание: вы можете передавать только ID, если хотите повысить безопасность
		Token: token,
	}

	return &response, nil
}

func (user *User) Login() (*AuthResponse, error) {
	var err error
	userFromDb := FetchUserByEmail(user.Email)

	if userFromDb.Email == "" {
		err = errors.New("User or password incorrect")
		return &AuthResponse{}, err
	}

	var isCheckedPassword = CheckPasswordHash(user.Password, userFromDb.Password)
	if !isCheckedPassword {
		err = errors.New("User or password incorrect")
		return &AuthResponse{}, err
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		return &AuthResponse{}, err
	}

	response := AuthResponse{
		User:  &userFromDb,
		Token: token,
	}

	return &response, nil
}

func (user *User) UpdateUser(id string) (*User, error) {
	if user.Password != "" {
		err := user.HashPassword()
		if err != nil {
			return &User{}, err
		}
	}

	err := Database.Model(&User{}).Where("id = ?", id).Updates(user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func GetUser(apiToken string) (*AuthResponse, error) {
	var authResponse AuthResponse
	err := Database.Where("token = ?", apiToken).First(&authResponse).Error

	if err != nil {
		return &AuthResponse{}, err
	}
	return &authResponse, nil
}
