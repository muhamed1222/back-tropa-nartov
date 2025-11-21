package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"tropa-nartov-backend/internal/logger"
	"tropa-nartov-backend/internal/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// StrapiClient клиент для синхронизации с Strapi
type StrapiClient struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
	db         *gorm.DB
}

// NewStrapiClient создает новый клиент Strapi
func NewStrapiClient(baseURL, apiToken string, db *gorm.DB) *StrapiClient {
	return &StrapiClient{
		baseURL:  baseURL,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		db: db,
	}
}

// StrapiPlace структура для Place из Strapi
type StrapiPlace struct {
	ID         int                    `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
}

// StrapiRoute структура для Route из Strapi
type StrapiRoute struct {
	ID         int                    `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
}

// SyncPlaces синхронизирует места из Strapi в Go DB
func (c *StrapiClient) SyncPlaces() error {
	logger.Info("Starting places synchronization from Strapi")
	
	// Запрос к Strapi API
	url := fmt.Sprintf("%s/api/places?pagination[limit]=100", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch places: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("strapi returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var result struct {
		Data []StrapiPlace `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	
	// Синхронизация каждого места
	synced := 0
	for _, strapiPlace := range result.Data {
		if err := c.syncPlace(strapiPlace); err != nil {
			logger.Error("Failed to sync place",
				zap.Int("strapi_id", strapiPlace.ID),
				zap.Error(err),
			)
			continue
		}
		synced++
	}
	
	logger.Info("Places synchronization completed",
		zap.Int("total", len(result.Data)),
		zap.Int("synced", synced),
	)
	
	return nil
}

func (c *StrapiClient) syncPlace(strapiPlace StrapiPlace) error {
	attrs := strapiPlace.Attributes
	
	// Проверяем существование по StrapiID
	var place models.Place
	err := c.db.Where("strapi_id = ?", strapiPlace.ID).First(&place).Error
	
	isNew := err == gorm.ErrRecordNotFound
	
	// Заполняем поля
	place.StrapiID = uint(strapiPlace.ID)
	place.Name = getString(attrs, "name")
	place.Description = getString(attrs, "description")
	place.Overview = getString(attrs, "overview")
	place.History = getString(attrs, "history")
	place.Address = getString(attrs, "address")
	place.Latitude = getFloat(attrs, "latitude")
	place.Longitude = getFloat(attrs, "longitude")
	place.Type = getString(attrs, "type")
	place.IsActive = true
	
	if isNew {
		return c.db.Create(&place).Error
	}
	return c.db.Save(&place).Error
}

// SyncRoutes синхронизирует маршруты из Strapi в Go DB
func (c *StrapiClient) SyncRoutes() error {
	logger.Info("Starting routes synchronization from Strapi")
	
	url := fmt.Sprintf("%s/api/routes?pagination[limit]=100", c.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch routes: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("strapi returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var result struct {
		Data []StrapiRoute `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	
	synced := 0
	for _, strapiRoute := range result.Data {
		if err := c.syncRoute(strapiRoute); err != nil {
			logger.Error("Failed to sync route",
				zap.Int("strapi_id", strapiRoute.ID),
				zap.Error(err),
			)
			continue
		}
		synced++
	}
	
	logger.Info("Routes synchronization completed",
		zap.Int("total", len(result.Data)),
		zap.Int("synced", synced),
	)
	
	return nil
}

func (c *StrapiClient) syncRoute(strapiRoute StrapiRoute) error {
	attrs := strapiRoute.Attributes
	
	var route models.Route
	err := c.db.Where("strapi_id = ?", strapiRoute.ID).First(&route).Error
	
	isNew := err == gorm.ErrRecordNotFound
	
	route.StrapiID = uint(strapiRoute.ID)
	route.Name = getString(attrs, "name")
	route.Description = getString(attrs, "description")
	route.Overview = getString(attrs, "overview")
	route.History = getString(attrs, "history")
	route.Distance = getFloat(attrs, "distance_km")
	route.Duration = getFloat(attrs, "duration_hours")
	route.IsActive = true
	
	if isNew {
		return c.db.Create(&route).Error
	}
	return c.db.Save(&route).Error
}

// WebhookHandler обрабатывает webhook от Strapi
func (c *StrapiClient) WebhookHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var webhook struct {
			Event string                 `json:"event"`
			Model string                 `json:"model"`
			Entry map[string]interface{} `json:"entry"`
		}
		
		if err := ctx.BindJSON(&webhook); err != nil {
			logger.Error("Failed to parse webhook", zap.Error(err))
			ctx.JSON(400, gin.H{"error": "invalid payload"})
			return
		}
		
		logger.Info("Received Strapi webhook",
			zap.String("event", webhook.Event),
			zap.String("model", webhook.Model),
		)
		
		// Обрабатываем webhook в зависимости от модели
		switch webhook.Model {
		case "place":
			// Синхронизировать конкретное место
			go c.SyncPlaces()
		case "route":
			// Синхронизировать конкретный маршрут
			go c.SyncRoutes()
		}
		
		ctx.JSON(200, gin.H{"status": "received"})
	}
}

// Helper functions
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getFloat(m map[string]interface{}, key string) float64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		}
	}
	return 0
}

