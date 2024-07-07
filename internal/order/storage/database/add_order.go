package storage

import (
	"WbTest/internal/order/model"
	"WbTest/internal/order/storage"
	"WbTest/internal/order/storage/database/dto"
	"context"
	"errors"
	"github.com/jackc/pgconn"
)

// func (s *OrderStorageDB) AddBanner(ctx context.Context, order model.Order) (*model.Order, error) {
func (s *OrderStorageDB) AddOrder(ctx context.Context, order model.Order) (*model.Order, error) {
	orderDB := dto.NewOrderDB(order)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				s.logger.Errorf("error in transaction rollback: %v", err)
				return
			}
		}
	}()

	// Insert into orders table
	_, err = tx.Exec(ctx,
		`INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		orderDB.OrderUID, orderDB.TrackNumber, orderDB.Entry, orderDB.Locale, orderDB.InternalSignature,
		orderDB.CustomerID, orderDB.DeliveryService, orderDB.ShardKey, orderDB.SmID, orderDB.DateCreated, orderDB.OofShard)
	if err != nil {
		return nil, err
	}

	// Insert into delivery table
	_, err = tx.Exec(ctx,
		`INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		orderDB.OrderUID, orderDB.Delivery.Name, orderDB.Delivery.Phone, orderDB.Delivery.Zip, orderDB.Delivery.City,
		orderDB.Delivery.Address, orderDB.Delivery.Region, orderDB.Delivery.Email)
	if err != nil {
		return nil, err
	}

	// Insert into payment table
	_, err = tx.Exec(ctx,
		`INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		orderDB.OrderUID, orderDB.Payment.Transaction, orderDB.Payment.RequestID, orderDB.Payment.Currency, orderDB.Payment.Provider,
		orderDB.Payment.Amount, orderDB.Payment.PaymentDt, orderDB.Payment.Bank, orderDB.Payment.DeliveryCost, orderDB.Payment.GoodsTotal, orderDB.Payment.CustomFee)
	if err != nil {
		return nil, err
	}

	// Insert into items table
	for _, item := range orderDB.Items {
		_, err = tx.Exec(ctx,
			`INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			orderDB.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return nil, storage.ErrDuplicateItem
			}
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
