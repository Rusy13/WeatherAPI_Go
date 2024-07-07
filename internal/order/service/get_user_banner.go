package service

import (
	"WbTest/internal/order/storage/database/dto"
	"context"
	"log"
)

func (s *OrderServiceApp) GetUserOrder(ctx context.Context, orderID string) (*dto.OrderFromCache, error) {
	// Попытка получить заказ из кэша
	log.Println("SERVICEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE")
	log.Println("ORDEEEEEEEEEEEEEEER = ", orderID)
	orderCache, err := s.storage.GetOrderFromCache(orderID)
	if err == nil {
		return orderCache, nil
	}

	// Если заказа нет в кэше, получить его из базы данных
	orderDB, err := s.storage.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Сохранить заказ в кэш для последующих запросов
	orderCache = dto.ConvertToOrderFromCache(*orderDB)
	err = s.storage.SaveOrderToCache(*orderCache, orderID)
	if err != nil {
		log.Panic("error saving order to cache: %v", err)
	}

	return orderCache, nil
}
