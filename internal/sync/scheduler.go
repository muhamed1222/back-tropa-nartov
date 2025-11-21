package sync

import (
	"tropa-nartov-backend/internal/logger"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Scheduler управляет периодической синхронизацией
type Scheduler struct {
	cron         *cron.Cron
	strapiClient *StrapiClient
}

// NewScheduler создает новый scheduler
func NewScheduler(strapiClient *StrapiClient) *Scheduler {
	return &Scheduler{
		cron:         cron.New(),
		strapiClient: strapiClient,
	}
}

// Start запускает периодическую синхронизацию
func (s *Scheduler) Start() error {
	// Синхронизация Places каждые 5 минут
	_, err := s.cron.AddFunc("*/5 * * * *", func() {
		logger.Info("Starting scheduled places synchronization")
		if err := s.strapiClient.SyncPlaces(); err != nil {
			logger.Error("Scheduled places sync failed", zap.Error(err))
		}
	})
	if err != nil {
		return err
	}

	// Синхронизация Routes каждые 5 минут (со смещением 2 минуты)
	_, err = s.cron.AddFunc("2,7,12,17,22,27,32,37,42,47,52,57 * * * *", func() {
		logger.Info("Starting scheduled routes synchronization")
		if err := s.strapiClient.SyncRoutes(); err != nil {
			logger.Error("Scheduled routes sync failed", zap.Error(err))
		}
	})
	if err != nil {
		return err
	}

	s.cron.Start()
	logger.Info("Sync scheduler started - running every 5 minutes")
	
	return nil
}

// Stop останавливает scheduler
func (s *Scheduler) Stop() {
	if s.cron != nil {
		s.cron.Stop()
		logger.Info("Sync scheduler stopped")
	}
}

