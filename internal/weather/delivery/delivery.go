package delivery

import (
	"WbTest/internal/weather/service"
	storage "WbTest/internal/weather/storage/database"
	"go.uber.org/zap"
)

type WeatherDelivery struct {
	service service.WeatherService
	storage *storage.WeatherStorageDB // Use a pointer to the storage instance
	logger  *zap.SugaredLogger
}

func New(service service.WeatherService, storage *storage.WeatherStorageDB, logger *zap.SugaredLogger) *WeatherDelivery {
	return &WeatherDelivery{
		service: service,
		storage: storage, // Assign the storage instance
		logger:  logger,
	}
}
