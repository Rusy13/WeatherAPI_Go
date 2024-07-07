package storage

import (
	"encoding/json"

	"WbTest/internal/order/storage/database/dto"
	"github.com/gomodule/redigo/redis"
)

func (s *OrderStorageDB) SaveOrderToCache(orderCache dto.OrderFromCache, orderID string) error {
	orderCacheJSON, err := json.Marshal(orderCache)
	if err != nil {
		s.logger.Errorf("error marshaling order to JSON: %v", err)
		return err
	}

	key := constructRedisKey(orderID)
	result, err := redis.String(s.redisConn.Do("SET", key, orderCacheJSON, "EX", s.cacheExpireTime))
	if err != nil || result != "OK" {
		s.logger.Errorf("error saving order to cache: %v", err)
	}
	return err
}

func constructRedisKey(orderID string) string {
	return "order_" + orderID
}
