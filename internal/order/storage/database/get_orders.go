package storage

import (
	"WbTest/internal/order/model"
	"WbTest/internal/order/storage/database/dto"
	"context"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
)

func (s *OrderStorageDB) GetAllOrders(ctx context.Context) ([]*model.Order, error) {
	var orders []*model.Order

	// Log the query
	s.logger.Info("Querying all orders from database")

	rows, err := s.db.Query(ctx,
		`SELECT 
			o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id, 
			o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
			d.name as delivery_name, d.phone as delivery_phone, d.zip as delivery_zip, 
			d.city as delivery_city, d.address as delivery_address, d.region as delivery_region, 
			d.email as delivery_email,
			p.transaction, p.request_id, p.currency, p.provider, p.amount, 
			p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM 
			orders o
		LEFT JOIN 
			delivery d ON o.order_uid = d.order_uid
		LEFT JOIN 
			payment p ON o.order_uid = p.order_uid`)

	if err != nil {
		s.logger.Errorf("Error querying orders: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var orderDB dto.OrderDB
		err := rows.Scan(
			&orderDB.OrderUID, &orderDB.TrackNumber, &orderDB.Entry, &orderDB.Locale, &orderDB.InternalSignature,
			&orderDB.CustomerID, &orderDB.DeliveryService, &orderDB.ShardKey, &orderDB.SmID, &orderDB.DateCreated, &orderDB.OofShard,
			&orderDB.Delivery.Name, &orderDB.Delivery.Phone, &orderDB.Delivery.Zip, &orderDB.Delivery.City,
			&orderDB.Delivery.Address, &orderDB.Delivery.Region, &orderDB.Delivery.Email,
			&orderDB.Payment.Transaction, &orderDB.Payment.RequestID, &orderDB.Payment.Currency, &orderDB.Payment.Provider,
			&orderDB.Payment.Amount, &orderDB.Payment.PaymentDt, &orderDB.Payment.Bank, &orderDB.Payment.DeliveryCost,
			&orderDB.Payment.GoodsTotal, &orderDB.Payment.CustomFee)
		if err != nil {
			s.logger.Errorf("Error scanning order row: %v", err)
			continue
		}

		order := dto.ConvertToOrder(orderDB)
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		s.logger.Errorf("Error with rows: %v", err)
		return nil, err
	}

	// Save orders to cache
	err = s.saveOrdersToCache(ctx, orders)
	if err != nil {
		s.logger.Errorf("Error saving orders to cache: %v", err)
	}

	return orders, nil
}

func (s *OrderStorageDB) saveOrdersToCache(ctx context.Context, orders []*model.Order) error {

	for _, order := range orders {
		orderCacheJSON, err := json.Marshal(order)
		if err != nil {
			s.logger.Errorf("Error marshaling order to JSON: %v", err)
			continue
		}

		key := constructRedisKey(order.OrderUID)
		result, err := redis.String(s.redisConn.Do("SET", key, orderCacheJSON, "EX", s.cacheExpireTime))
		if err != nil || result != "OK" {
			s.logger.Errorf("error saving order to cache: %v", err)
		}
		return err
	}

	return nil
}
