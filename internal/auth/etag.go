package auth

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// ETagMiddleware создает middleware для поддержки ETag кеширования
// Генерирует ETag на основе содержимого ответа и проверяет If-None-Match заголовок
func ETagMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пропускаем запрос, если это не GET/HEAD
		if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			c.Next()
			return
		}

		// Перехватываем ответ
		writer := &responseWriter{body: make([]byte, 0), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		// Проверяем статус ответа - обрабатываем только успешные ответы
		if writer.Status() < 200 || writer.Status() >= 300 {
			// Отправляем ответ как есть без ETag
			c.Writer = writer.ResponseWriter
			c.Writer.WriteHeader(writer.Status())
			c.Writer.Write(writer.body)
			return
		}

		// Генерируем ETag на основе содержимого
		hash := md5.Sum(writer.body)
		etag := fmt.Sprintf(`"%s"`, hex.EncodeToString(hash[:]))

		// Проверяем If-None-Match заголовок
		ifNoneMatch := c.GetHeader("If-None-Match")
		if ifNoneMatch != "" {
			// Сравниваем ETag (может быть список через запятую)
			ifNoneMatchValues := strings.Split(ifNoneMatch, ",")
			for _, value := range ifNoneMatchValues {
				value = strings.TrimSpace(value)
				if value == etag {
					// Контент не изменился - возвращаем 304 Not Modified
					c.Writer = writer.ResponseWriter
					c.Writer.Header().Set("ETag", etag)
					c.Writer.WriteHeader(304)
					return
				}
			}
		}

		// Контент изменился или запрос без If-None-Match - отправляем полный ответ
		c.Writer = writer.ResponseWriter
		c.Writer.Header().Set("ETag", etag)
		c.Writer.Header().Set("Cache-Control", "public, max-age=300") // 5 минут кеширования
		c.Writer.WriteHeader(writer.Status())
		c.Writer.Write(writer.body)
	}
}

// responseWriter перехватывает запись ответа для генерации ETag
type responseWriter struct {
	gin.ResponseWriter
	body   []byte
	status int
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return len(b), nil
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body = append(w.body, []byte(s)...)
	return len(s), nil
}

func (w *responseWriter) Status() int {
	if w.status == 0 {
		return 200
	}
	return w.status
}

// Реализуем io.Writer для совместимости
func (w *responseWriter) WriteTo(dst io.Writer) (int64, error) {
	n, err := dst.Write(w.body)
	return int64(n), err
}

