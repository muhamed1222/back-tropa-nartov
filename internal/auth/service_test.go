package auth

import (
	"testing"
	"time"
	"tropa-nartov-backend/internal/config"
	"tropa-nartov-backend/internal/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}
	
	// Migrate test schema
	db.AutoMigrate(&models.User{})
	
	return db
}

func TestRegister(t *testing.T) {
	db := setupTestDB()
	cfg := &config.Config{
		JWTSecret:    "test-secret",
		JWTExpiresIn: 24,
	}
	service := NewService(db, cfg)

	// Test успешной регистрации
	user, err := service.Register("Test User", "test@example.com", "password123")
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotEmpty(t, user.PasswordHash)

	// Test duplicate email
	_, err = service.Register("Another User", "test@example.com", "password456")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "уже существует")
}

func TestLogin(t *testing.T) {
	db := setupTestDB()
	cfg := &config.Config{
		JWTSecret:    "test-secret",
		JWTExpiresIn: 24,
	}
	service := NewService(db, cfg)

	// Создаем тестового пользователя
	_, err := service.Register("Test User", "test@example.com", "password123")
	assert.NoError(t, err)

	// Test успешного логина
	token, err := service.Login("test@example.com", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test неверного пароля
	_, err = service.Login("test@example.com", "wrongpassword")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "неверный пароль")

	// Test несуществующего пользователя
	_, err = service.Login("nonexistent@example.com", "password123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "не найден")
}

func TestGetUserProfile(t *testing.T) {
	db := setupTestDB()
	cfg := &config.Config{
		JWTSecret:    "test-secret",
		JWTExpiresIn: 24,
	}
	service := NewService(db, cfg)

	// Создаем пользователя
	user, err := service.Register("Test User", "test@example.com", "password123")
	assert.NoError(t, err)

	// Получаем профиль
	profile, err := service.GetUserProfile(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, profile.Email)

	// Test несуществующего пользователя
	_, err = service.GetUserProfile(999)
	assert.Error(t, err)
}

func TestUpdateUserProfile(t *testing.T) {
	db := setupTestDB()
	cfg := &config.Config{
		JWTSecret:    "test-secret",
		JWTExpiresIn: 24,
	}
	service := NewService(db, cfg)

	// Создаем пользователя
	user, err := service.Register("Test User", "test@example.com", "password123")
	assert.NoError(t, err)

	// Обновляем профиль
	updated, err := service.UpdateUserProfile(user.ID, "Updated Name", "Updated FirstName", "test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, "Updated FirstName", updated.FirstName)
}

func TestChangePassword(t *testing.T) {
	db := setupTestDB()
	cfg := &config.Config{
		JWTSecret:    "test-secret",
		JWTExpiresIn: 24,
	}
	service := NewService(db, cfg)

	// Создаем пользователя
	user, err := service.Register("Test User", "test@example.com", "oldpassword123")
	assert.NoError(t, err)

	// Меняем пароль
	err = service.ChangePassword(user.ID, "oldpassword123", "newpassword456")
	assert.NoError(t, err)

	// Проверяем что можем войти с новым паролем
	token, err := service.Login("test@example.com", "newpassword456")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Проверяем что старый пароль не работает
	_, err = service.Login("test@example.com", "oldpassword123")
	assert.Error(t, err)
}

func TestForgotPassword(t *testing.T) {
	db := setupTestDB()
	cfg := &config.Config{
		JWTSecret:    "test-secret",
		JWTExpiresIn: 24,
		Environment:  "development",
	}
	service := NewService(db, cfg)

	// Создаем пользователя
	user, err := service.Register("Test User", "test@example.com", "password123")
	assert.NoError(t, err)

	// Запрашиваем сброс пароля
	err = service.ForgotPassword("test@example.com")
	assert.NoError(t, err)

	// Проверяем что код создан
	var updatedUser models.User
	db.First(&updatedUser, user.ID)
	assert.NotEmpty(t, updatedUser.ResetToken)
	assert.False(t, updatedUser.ResetTokenExpires.IsZero())
}

func TestResetPassword(t *testing.T) {
	db := setupTestDB()
	cfg := &config.Config{
		JWTSecret:    "test-secret",
		JWTExpiresIn: 24,
		Environment:  "development",
	}
	service := NewService(db, cfg)

	// Создаем пользователя
	user, err := service.Register("Test User", "test@example.com", "oldpassword")
	assert.NoError(t, err)

	// Создаем reset token вручную для теста
	user.ResetToken = "123456"
	user.ResetTokenExpires = time.Now().Add(15 * time.Minute)
	db.Save(&user)

	// Сбрасываем пароль
	err = service.ResetPassword("123456", "newpassword123")
	assert.NoError(t, err)

	// Проверяем что можем войти с новым паролем
	token, err := service.Login("test@example.com", "newpassword123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

