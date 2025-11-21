package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"
	"tropa-nartov-backend/internal/config"
	"tropa-nartov-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	cfg       *config.Config
	jwtSecret string
}

func NewService(db *gorm.DB, cfg *config.Config) *Service {
	jwtSecret := cfg.JWTSecret
	if jwtSecret == "" {
		panic("FATAL: JWT_SECRET_KEY must be set in environment variables")
	}

	return &Service{
		db:        db,
		cfg:       cfg,
		jwtSecret: jwtSecret,
	}
}

// GetUserProfile возвращает данные пользователя по ID
func (s *Service) GetUserProfile(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ? AND is_active = true", userID).First(&user).Error; err != nil {
		return nil, errors.New("пользователь не найден")
	}
	return &user, nil
}

// Register создает нового пользователя
func (s *Service) Register(firstName, email, password string) (*models.User, error) {
	// Проверяем, существует ли пользователь
	var existingUser models.User
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("ошибка хеширования пароля")
	}

	// Создаем пользователя
	user := &models.User{
		Name:         firstName, // Сохраняем как name
		FirstName:    "",        // Фамилия пустая
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
		IsActive:     true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, errors.New("ошибка создания пользователя: " + err.Error())
	}

	return user, nil
}

// UpdateUserProfile обновляет данные пользователя
func (s *Service) UpdateUserProfile(userID uint, name, firstName, email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ? AND is_active = true", userID).First(&user).Error; err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// Проверяем email на уникальность, если он изменился
	if user.Email != email {
		var existingUser models.User
		if err := s.db.Where("email = ? AND id != ?", email, userID).First(&existingUser).Error; err == nil {
			return nil, errors.New("пользователь с таким email уже существует")
		}
		user.Email = email
	}

	user.Name = name
	user.FirstName = firstName // Фамилия
	user.UpdatedAt = time.Now()

	if err := s.db.Save(&user).Error; err != nil {
		return nil, errors.New("ошибка обновления профиля")
	}

	return &user, nil
}

// Login аутентифицирует пользователя
func (s *Service) Login(email, password string) (string, error) {
	var user models.User
	if err := s.db.Where("email = ? AND is_active = true", email).First(&user).Error; err != nil {
		return "", errors.New("пользователь не найден")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("неверный пароль")
	}

	// Создаем JWT токен с правильным форматом
	// Используем JWTExpiresIn из конфигурации (в часах)
	expiresInHours := s.cfg.JWTExpiresIn
	if expiresInHours <= 0 {
		expiresInHours = 24 // Значение по умолчанию если не задано
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Duration(expiresInHours) * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})

	// Подписываем токен
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", errors.New("ошибка создания токена")
	}

	return tokenString, nil
}

// generateResetCode генерирует короткий код для сброса пароля (6 цифр)
func (s *Service) generateResetCode() (string, error) {
	// Генерируем 6-значный код
	bytes := make([]byte, 3)
	if _, err := rand.Read(bytes); err != nil {
		return "", errors.New("ошибка генерации кода")
	}
	// Преобразуем в число от 100000 до 999999
	code := 100000 + (int(bytes[0])<<16|int(bytes[1])<<8|int(bytes[2]))%900000
	return fmt.Sprintf("%06d", code), nil
}

// ForgotPassword инициирует процесс сброса пароля
func (s *Service) ForgotPassword(email string) error {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		// Не раскрываем информацию о том, существует ли пользователь
		// Всегда возвращаем успех для безопасности
		return nil
	}

	// Генерируем код сброса пароля (6 цифр)
	code, err := s.generateResetCode()
	if err != nil {
		return errors.New("ошибка генерации кода сброса пароля")
	}

	// Сохраняем код и время истечения (15 минут для кода)
	user.ResetToken = code
	user.ResetTokenExpires = time.Now().Add(15 * time.Minute)

	if err := s.db.Save(&user).Error; err != nil {
		return errors.New("ошибка сохранения кода сброса пароля")
	}

	// Отправляем email с кодом
	err = SendPasswordResetEmail(s.cfg, user.Email, code)
	if err != nil {
		// В dev режиме не возвращаем ошибку клиенту, только логируем
		if s.cfg.Debug || s.cfg.Environment == "development" {
			// log.Printf("⚠️  Ошибка отправки email: %v", err)
			// log.Printf("🔐 [DEV MODE] Код сброса пароля для %s: %s", user.Email, code)
			// log.Printf("🔐 [DEV MODE] Используйте этот код для сброса пароля")
			// В dev режиме не возвращаем ошибку - код сохранен и доступен в логах
			return nil
		}
		// В production возвращаем ошибку, если email не отправлен
		return fmt.Errorf("не удалось отправить код на email: %w", err)
	}

	return nil
}

// VerifyResetCode проверяет код восстановления
func (s *Service) VerifyResetCode(code string) error {
	if code == "" {
		return errors.New("код не может быть пустым")
	}

	var user models.User
	if err := s.db.Where("reset_token = ?", code).First(&user).Error; err != nil {
		return errors.New("неверный или истёкший код")
	}

	// Проверяем, не истек ли код
	if user.ResetTokenExpires.Before(time.Now()) {
		return errors.New("код истёк. Запросите новый код")
	}

	return nil
}

// ResetPassword сбрасывает пароль пользователя
func (s *Service) ResetPassword(code, newPassword string) error {
	if code == "" {
		return errors.New("код не может быть пустым")
	}

	if newPassword == "" || len(newPassword) < 8 {
		return errors.New("пароль должен содержать минимум 8 символов")
	}

	// Находим пользователя по коду
	var user models.User
	if err := s.db.Where("reset_token = ?", code).First(&user).Error; err != nil {
		return errors.New("неверный или истёкший код")
	}

	// Проверяем, не истек ли код
	if user.ResetTokenExpires.Before(time.Now()) {
		return errors.New("код истёк. Пожалуйста, запросите новый код для сброса пароля")
	}

	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("ошибка хеширования пароля")
	}

	// Обновляем пароль и очищаем токен сброса
	user.PasswordHash = string(hashedPassword)
	user.ResetToken = ""
	user.ResetTokenExpires = time.Time{}
	user.UpdatedAt = time.Now()

	if err := s.db.Save(&user).Error; err != nil {
		return errors.New("ошибка обновления пароля")
	}

	// fmt.Printf("✅ Пароль успешно сброшен для пользователя: %s\n", user.Email)
	return nil
}

// DeleteUser полностью удаляет пользователя и все связанные данные
func (s *Service) DeleteUser(userID int) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// Удаляем связанные данные в правильном порядке
	if err := tx.Where("user_id = ?", userID).Delete(&models.FavoritePlace{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&models.FavoriteRoute{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&models.PassedRoute{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&models.Review{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// ChangePassword меняет пароль пользователя
func (s *Service) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// Валидация нового пароля
	if newPassword == "" || len(newPassword) < 8 {
		return errors.New("пароль должен содержать минимум 8 символов")
	}

	// Проверяем, что новый пароль отличается от старого
	if oldPassword == newPassword {
		return errors.New("новый пароль должен отличаться от старого")
	}

	var user models.User
	if err := s.db.Where("id = ? AND is_active = true", userID).First(&user).Error; err != nil {
		return errors.New("пользователь не найден")
	}

	// Проверяем старый пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("неверный старый пароль")
	}

	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("ошибка хеширования пароля")
	}

	// Обновляем пароль
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.db.Save(&user).Error; err != nil {
		return errors.New("ошибка обновления пароля")
	}

	return nil
}
