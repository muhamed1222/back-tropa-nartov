package validation

import (
	"html"
	"regexp"
	"strings"
)

// SanitizeHTML очищает HTML теги из строки
func SanitizeHTML(input string) string {
	// Убираем HTML теги
	input = html.EscapeString(input)
	// Убираем лишние пробелы
	input = strings.TrimSpace(input)
	return input
}

// SanitizeText очищает текст от опасных символов
func SanitizeText(input string) string {
	input = strings.TrimSpace(input)
	// Убираем control characters
	re := regexp.MustCompile(`[\x00-\x1F\x7F]`)
	input = re.ReplaceAllString(input, "")
	return input
}

// ValidateEmail проверяет формат email
func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// ValidatePhone проверяет формат телефона (упрощенно)
func ValidatePhone(phone string) bool {
	// Убираем все кроме цифр и +
	re := regexp.MustCompile(`^[\+]?[0-9\s\-\(\)]{10,15}$`)
	return re.MatchString(phone)
}

// ValidatePassword проверяет сложность пароля
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Пароль должен содержать минимум 8 символов"
	}
	if len(password) > 128 {
		return false, "Пароль слишком длинный (максимум 128 символов)"
	}
	
	// Проверяем наличие букв и цифр
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	
	if !hasLetter || !hasNumber {
		return false, "Пароль должен содержать буквы и цифры"
	}
	
	return true, ""
}

// SanitizeFilename очищает имя файла
func SanitizeFilename(filename string) string {
	// Убираем опасные символы из имени файла
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	filename = re.ReplaceAllString(filename, "_")
	
	// Ограничиваем длину
	if len(filename) > 255 {
		filename = filename[:255]
	}
	
	return filename
}

