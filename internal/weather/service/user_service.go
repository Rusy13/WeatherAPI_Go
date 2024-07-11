package service

import (
	"WbTest/internal/weather/model"
	"WbTest/internal/weather/storage"
	"context"
)

// UserService определяет методы для работы с пользователями и их избранными городами.
type UserService interface {
	RegisterUser(ctx context.Context, user *model.User) error
	LoginUser(ctx context.Context, user *model.User) error
	AddFavoriteCity(ctx context.Context, favCity *model.FavoriteCity) error
}

// UserServiceImpl реализует интерфейс UserService.
type UserServiceImpl struct {
	storage storage.UserStorage
}

// NewUserService создает новый экземпляр UserServiceImpl.
func NewUserService(storage storage.UserStorage) *UserServiceImpl {
	return &UserServiceImpl{
		storage: storage,
	}
}

func (s *UserServiceImpl) RegisterUser(ctx context.Context, user *model.User) error {
	return s.storage.RegisterUser(ctx, user)
}

func (s *UserServiceImpl) LoginUser(ctx context.Context, user *model.User) error {
	return s.storage.LoginUser(ctx, user)
}

func (s *UserServiceImpl) AddFavoriteCity(ctx context.Context, favCity *model.FavoriteCity) error {
	return s.storage.AddFavoriteCity(ctx, favCity)
}
