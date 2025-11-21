package middleware

import (
	"fmt"
	"net/http"
	"time"
	"tropa-nartov-backend/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"go.uber.org/zap"
)

// RateLimitMiddleware создает middleware для rate limiting по IP адресу
func RateLimitMiddleware(rate limiter.Rate) gin.HandlerFunc {
	// Создаем in-memory store
	store := memory.NewStore()
	
	// Создаем rate limiter instance
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Получаем IP адрес клиента
		ip := c.ClientIP()
		
		// Проверяем лимит
		context, err := instance.Get(c, ip)
		if err != nil {
			logger.Error("Rate limiter error", 
				zap.Error(err),
				zap.String("ip", ip),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка проверки rate limit",
			})
			c.Abort()
			return
		}

		// Добавляем заголовки с информацией о лимитах
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", context.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", context.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", context.Reset))

		// Если лимит превышен
		if context.Reached {
			logger.Warn("Rate limit exceeded",
				zap.String("ip", ip),
				zap.String("path", c.Request.URL.Path),
			)
			
			// Вычисляем сколько секунд до сброса
			retryAfter := int64(0)
			if context.Reset > time.Now().Unix() {
				retryAfter = context.Reset - time.Now().Unix()
			}
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Слишком много запросов",
				"message":     "Вы превысили лимит запросов. Попробуйте позже.",
				"retry_after": retryAfter,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GlobalRateLimit - глобальный лимит для всех эндпоинтов (100 запросов в минуту)
func GlobalRateLimit() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}
	return RateLimitMiddleware(rate)
}

// AuthRateLimit - лимит для auth эндпоинтов (10 запросов в минуту)
func AuthRateLimit() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  10,
	}
	return RateLimitMiddleware(rate)
}

// APIRateLimit - лимит для API эндпоинтов (30 запросов в минуту)
func APIRateLimit() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  30,
	}
	return RateLimitMiddleware(rate)
}

