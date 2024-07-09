package storage

import (
	"WbTest/internal/order/model"
	"context"
)

// UserStorage определяет методы для работы с данными пользователей.
type UserStorage interface {
	RegisterUser(ctx context.Context, user *model.User) error
	LoginUser(ctx context.Context, user *model.User) error
	AddFavoriteCity(ctx context.Context, favCity *model.FavoriteCity) error
}
