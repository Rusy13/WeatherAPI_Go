package dto

import (
	"WbTest/internal/order/model"
	"fmt"
	"github.com/asaskevich/govalidator"
	"time"
)

type AddOrderDTO struct {
	OrderUID           string    `json:"order_uid" valid:"required"`
	TrackNumber        string    `json:"track_number" valid:"required"`
	Entry              string    `json:"entry" valid:"required"`
	DeliveryName       string    `json:"delivery_name" valid:"required"`
	DeliveryPhone      string    `json:"delivery_phone" valid:"required"`
	DeliveryZip        string    `json:"delivery_zip" valid:"required"`
	DeliveryCity       string    `json:"delivery_city" valid:"required"`
	DeliveryAddress    string    `json:"delivery_address" valid:"required"`
	DeliveryRegion     string    `json:"delivery_region" valid:"required"`
	DeliveryEmail      string    `json:"delivery_email" valid:"required,email"`
	PaymentTransaction string    `json:"payment_transaction" valid:"required"`
	RequestID          string    `json:"request_id"`
	Currency           string    `json:"currency" valid:"required"`
	Provider           string    `json:"provider" valid:"required"`
	Amount             int       `json:"amount" valid:"required"`
	PaymentDt          int       `json:"payment_dt" valid:"required"`
	Bank               string    `json:"bank" valid:"required"`
	DeliveryCost       int       `json:"delivery_cost" valid:"required"`
	GoodsTotal         int       `json:"goods_total" valid:"required"`
	CustomFee          int       `json:"custom_fee" valid:"required"`
	Items              []ItemDTO `json:"items" valid:"required"`
	Locale             string    `json:"locale"`
	InternalSignature  string    `json:"internal_signature"`
	CustomerID         string    `json:"customer_id"`
	DeliveryService    string    `json:"delivery_service"`
	Shardkey           string    `json:"shardkey"`
	SmID               int       `json:"sm_id"`
	DateCreated        string    `json:"date_created" valid:"required"`
	OofShard           string    `json:"oof_shard"`
}

type ItemDTO struct {
	ChrtID      int    `json:"chrt_id" valid:"required"`
	TrackNumber string `json:"track_number" valid:"required"`
	Price       int    `json:"price" valid:"required"`
	Rid         string `json:"rid" valid:"required"`
	Name        string `json:"name" valid:"required"`
	Sale        int    `json:"sale" valid:"required"`
	Size        string `json:"size" valid:"required"`
	TotalPrice  int    `json:"total_price" valid:"required"`
	NmID        int    `json:"nm_id" valid:"required"`
	Brand       string `json:"brand" valid:"required"`
	Status      int    `json:"status" valid:"required"`
}

func (a *AddOrderDTO) Validate() error {
	_, err := govalidator.ValidateStruct(a)
	if err != nil {
		fmt.Printf("Validation errors: %v\n", err)
	}
	return err
}

func ConvertToOrder(b AddOrderDTO) model.Order {
	return model.Order{
		OrderUID:    b.OrderUID,
		TrackNumber: b.TrackNumber,
		Entry:       b.Entry,
		Delivery: model.Delivery{
			Name:    b.DeliveryName,
			Phone:   b.DeliveryPhone,
			Zip:     b.DeliveryZip,
			City:    b.DeliveryCity,
			Address: b.DeliveryAddress,
			Region:  b.DeliveryRegion,
			Email:   b.DeliveryEmail,
		},
		Payment: model.Payment{
			Transaction:  b.PaymentTransaction,
			RequestID:    b.RequestID,
			Currency:     b.Currency,
			Provider:     b.Provider,
			Amount:       b.Amount,
			PaymentDt:    b.PaymentDt,
			Bank:         b.Bank,
			DeliveryCost: b.DeliveryCost,
			GoodsTotal:   b.GoodsTotal,
			CustomFee:    b.CustomFee,
		},
		Items: func(itemsDTO []ItemDTO) []model.Item {
			items := make([]model.Item, len(itemsDTO))
			for i, itemDTO := range itemsDTO {
				items[i] = model.Item{
					ChrtID:      itemDTO.ChrtID,
					TrackNumber: itemDTO.TrackNumber,
					Price:       itemDTO.Price,
					Rid:         itemDTO.Rid,
					Name:        itemDTO.Name,
					Sale:        itemDTO.Sale,
					Size:        itemDTO.Size,
					TotalPrice:  itemDTO.TotalPrice,
					NmID:        itemDTO.NmID,
					Brand:       itemDTO.Brand,
					Status:      itemDTO.Status,
				}
			}
			return items
		}(b.Items),
		Locale:            b.Locale,
		InternalSignature: b.InternalSignature,
		CustomerID:        b.CustomerID,
		DeliveryService:   b.DeliveryService,
		Shardkey:          b.Shardkey,
		SmID:              b.SmID,
		DateCreated: func(dateCreatedStr string) time.Time {
			dateCreated, _ := time.Parse(time.RFC3339, dateCreatedStr)
			return dateCreated
		}(b.DateCreated),
		OofShard: b.OofShard,
	}
}
