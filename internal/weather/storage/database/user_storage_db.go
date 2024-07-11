package storage

import (
	"WbTest/internal/infrastructure/database/postgres/database"
	"WbTest/internal/weather/model"
	"context"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserStorageDB struct {
	db     database.Database
	logger *zap.SugaredLogger
}

func NewUserStorageDB(db database.Database, logger *zap.SugaredLogger) *UserStorageDB {
	return &UserStorageDB{db: db, logger: logger}
}

func (s *UserStorageDB) RegisterUser(ctx context.Context, user *model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	err = s.db.QueryRow(ctx, query, user.Username, hashedPassword).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

func (s *UserStorageDB) LoginUser(ctx context.Context, user *model.User) error {
	var hashedPassword string
	query := `SELECT id, password FROM users WHERE username = $1`
	err := s.db.QueryRow(ctx, query, user.Username).Scan(&user.ID, &hashedPassword)
	if err != nil {
		return fmt.Errorf("invalid username or password: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		return fmt.Errorf("invalid username or password: %w", err)
	}

	return nil
}

func (s *UserStorageDB) AddFavoriteCity(ctx context.Context, favCity *model.FavoriteCity) error {
	// Check if the city exists in the database
	var cityExists bool
	checkCityQuery := `SELECT EXISTS(SELECT 1 FROM cities WHERE name = $1)`
	err := s.db.QueryRow(ctx, checkCityQuery, favCity.CityName).Scan(&cityExists)
	if err != nil {
		return fmt.Errorf("failed to check if city exists: %w", err)
	}

	if !cityExists {
		return fmt.Errorf("city %s does not exist", favCity.CityName)
	}

	// Add the favorite city
	query := `INSERT INTO favorite_cities (user_id, city_name) VALUES ($1, $2)`
	_, err = s.db.Exec(ctx, query, favCity.UserID, favCity.CityName)
	if err != nil {
		return fmt.Errorf("failed to add favorite city: %w", err)
	}

	return nil
}
