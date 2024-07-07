package service

import (
	"WbTest/internal/order/model"
	"WbTest/internal/order/storage/database/dto"
	"context"
	"log"
	"time"
)

func (s *OrderServiceApp) AddOrder(ctx context.Context, order model.Order) (*model.Order, error) {
	creationTime := time.Now()
	order.DateCreated = creationTime

	addedOrder, err := s.storage.AddOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// Сохранение заказа в кэш
	orderCache := dto.ConvertToOrderFromCache(*addedOrder)
	err = s.storage.SaveOrderToCache(*orderCache, addedOrder.OrderUID)
	if err != nil {
		log.Panic("error saving order to cache: %v", err)
	}

	return addedOrder, nil
}
