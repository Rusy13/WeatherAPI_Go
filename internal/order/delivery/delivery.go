package delivery

import (
	"WbTest/internal/order/service"
	"go.uber.org/zap"
)

type OrderDelivery struct {
	service service.OrderService
	logger  *zap.SugaredLogger
}

func New(service service.OrderService, logger *zap.SugaredLogger) *OrderDelivery {
	return &OrderDelivery{
		service: service,
		logger:  logger,
	}
}
