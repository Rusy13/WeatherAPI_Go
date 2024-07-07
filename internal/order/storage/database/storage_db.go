package storage

import (
	"WbTest/internal/infrastructure/database/postgres/database"
	"go.uber.org/zap"

	"github.com/gomodule/redigo/redis"
)

const expireTime = 15

type OrderStorageDB struct {
	db              database.Database
	redisConn       redis.Conn
	cacheExpireTime int
	logger          *zap.SugaredLogger
}

func New(db database.Database, redisConn redis.Conn, logger *zap.SugaredLogger) *OrderStorageDB {
	return &OrderStorageDB{
		db:              db,
		logger:          logger,
		redisConn:       redisConn,
		cacheExpireTime: expireTime,
	}
}
