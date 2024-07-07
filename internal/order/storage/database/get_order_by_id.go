package storage

import (
	"WbTest/internal/order/model"
	"WbTest/internal/order/storage"
	"WbTest/internal/order/storage/database/dto"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
)

func (s *OrderStorageDB) GetOrderByID(ctx context.Context, orderUID string) (*model.Order, error) {
	var orderDB dto.OrderDB

	// Log the query and parameter
	s.logger.Infof("Querying order details for orderUID: %s", orderUID)

	err := s.db.QueryRow(ctx,
		`SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
         FROM orders WHERE order_uid = $1`, orderUID).Scan(
		&orderDB.OrderUID, &orderDB.TrackNumber, &orderDB.Entry, &orderDB.Locale, &orderDB.InternalSignature,
		&orderDB.CustomerID, &orderDB.DeliveryService, &orderDB.ShardKey, &orderDB.SmID, &orderDB.DateCreated, &orderDB.OofShard)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Infof("Order not found in DB: %s", orderUID)
			return nil, storage.ErrOrderNotFound
		}
		s.logger.Errorf("Error querying order details: %v", err)
		return nil, err
	}

	s.logger.Infof("Order details fetched: %+v", orderDB)

	// Fetch delivery details
	s.logger.Infof("Querying delivery details for orderUID: %s", orderUID)
	err = s.db.QueryRow(ctx,
		`SELECT name, phone, zip, city, address, region, email 
         FROM delivery WHERE order_uid = $1`, orderUID).Scan(
		&orderDB.Delivery.Name, &orderDB.Delivery.Phone, &orderDB.Delivery.Zip, &orderDB.Delivery.City,
		&orderDB.Delivery.Address, &orderDB.Delivery.Region, &orderDB.Delivery.Email)
	if err != nil {
		s.logger.Errorf("Error querying delivery details: %v", err)
		return nil, err
	}

	s.logger.Infof("Delivery details fetched: %+v", orderDB.Delivery)

	// Fetch payment details
	s.logger.Infof("Querying payment details for orderUID: %s", orderUID)
	err = s.db.QueryRow(ctx,
		`SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
         FROM payment WHERE order_uid = $1`, orderUID).Scan(
		&orderDB.Payment.Transaction, &orderDB.Payment.RequestID, &orderDB.Payment.Currency, &orderDB.Payment.Provider,
		&orderDB.Payment.Amount, &orderDB.Payment.PaymentDt, &orderDB.Payment.Bank, &orderDB.Payment.DeliveryCost,
		&orderDB.Payment.GoodsTotal, &orderDB.Payment.CustomFee)
	if err != nil {
		s.logger.Errorf("Error querying payment details: %v", err)
		return nil, err
	}

	s.logger.Infof("Payment details fetched: %+v", orderDB.Payment)

	// Fetch items details
	s.logger.Infof("Querying items for orderUID: %s", orderUID)
	rows, err := s.db.Query(ctx,
		`SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
         FROM items WHERE order_uid = $1`, orderUID)
	if err != nil {
		s.logger.Errorf("Error querying items: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item dto.ItemDB
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			s.logger.Errorf("Error scanning item: %v", err)
			return nil, err
		}
		orderDB.Items = append(orderDB.Items, item)
	}

	if err = rows.Err(); err != nil {
		s.logger.Errorf("Error with rows: %v", err)
		return nil, err
	}

	s.logger.Infof("Items fetched: %+v", orderDB.Items)

	order := dto.ConvertToOrder(orderDB)
	return &order, nil
}
