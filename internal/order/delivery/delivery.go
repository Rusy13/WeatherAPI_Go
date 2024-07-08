package delivery

import (
	"WbTest/internal/order/service"
	storage "WbTest/internal/order/storage/database"
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
