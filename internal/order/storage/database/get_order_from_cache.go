package storage

import (
	"WbTest/internal/order/storage/database/dto"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"log"
)

func (s *OrderStorageDB) GetOrderFromCache(orderID string) (*dto.OrderFromCache, error) {
	key := constructRedisKey(orderID)
	orderFromRedis, err := redis.String(s.redisConn.Do("GET", key))
	if err != nil {
		s.logger.Errorf("redis error: %v", err)
		log.Printf("Redis error: %v", err)
		return nil, err
	}

	var orderCache dto.OrderFromCache
	err = json.Unmarshal([]byte(orderFromRedis), &orderCache)
	if err != nil {
		s.logger.Errorf("unmarshal error: %v", err)
		log.Printf("Unmarshal error: %v", err)
		return nil, err
	}

	log.Printf("Order fetched from Redis cache: %+v", orderCache)
	log.Println("DATA FROM CACHEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEES")
	return &orderCache, nil
}
