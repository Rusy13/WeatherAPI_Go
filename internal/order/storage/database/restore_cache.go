package storage

import (
	"WbTest/internal/order/storage/database/dto"
	"context"
	"fmt"
	"log"
)

func (s *OrderStorageDB) RestoreCacheFromDB() error {
	// Здесь выполняется логика восстановления кэша из базы данных

	// 1. Получаем все заказы из базы данных
	ctx := context.Background()
	orders, err := s.GetAllOrders(ctx)
	if err != nil {
		return fmt.Errorf("error fetching orders from database: %v", err)
	}

	// 2. Проходимся по каждому заказу и записываем его в кэш
	for _, order := range orders {
		// Преобразуем данные заказа в формат кэша
		orderCache := dto.ConvertToOrderCache(*order)

		// Сохраняем данные заказа в кэше
		err := s.SaveOrderToCache(*orderCache, order.OrderUID)
		if err != nil {
			// В случае ошибки сохранения в кэше, можно просто логировать её
			// или вернуть ошибку, чтобы обработать её выше
			log.Printf("error saving order %s to cache: %v", order.OrderUID, err)
		}
	}

	return nil
}
