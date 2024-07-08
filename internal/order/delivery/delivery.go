package delivery

import (
	"WbTest/internal/order/service"
	"go.uber.org/zap"
)

type WeatherDelivery struct {
	service service.WeatherService
	logger  *zap.SugaredLogger
}

func New(service service.WeatherService, logger *zap.SugaredLogger) *WeatherDelivery {
	return &WeatherDelivery{
		service: service,
		logger:  logger,
	}
}
