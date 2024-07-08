package storage

import (
	"WbTest/internal/infrastructure/database/postgres/database"
	"WbTest/internal/order/model"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

const expireTime = 15

type WeatherStorageDB struct {
	db              database.Database
	redisConn       redis.Conn
	cacheExpireTime int
	logger          *zap.SugaredLogger
}

func New(db database.Database, redisConn redis.Conn, logger *zap.SugaredLogger) *WeatherStorageDB {
	return &WeatherStorageDB{
		db:              db,
		redisConn:       redisConn,
		cacheExpireTime: expireTime,
		logger:          logger,
	}
}

func (s *WeatherStorageDB) GetCountryForCity(ctx context.Context, city string) (string, error) {
	query := `SELECT country FROM cities WHERE name = $1`
	row := s.db.QueryRow(ctx, query, city)

	var country string
	err := row.Scan(&country)
	if err != nil {
		return "", fmt.Errorf("failed to get country for city %s: %w", city, err)
	}

	return country, nil
}

func (s *WeatherStorageDB) SaveWeather(ctx context.Context, city string, weather model.Weather) error {
	query := `
        INSERT INTO weather (city_name, date, temp, data)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (city_name, date) DO UPDATE 
        SET temp = EXCLUDED.temp,
            data = EXCLUDED.data;
    `

	rawData, err := json.Marshal(weather)
	if err != nil {
		return fmt.Errorf("failed to marshal weather data: %w", err)
	}

	_, err = s.db.Exec(ctx, query, city, weather.DateTime, weather.Temperature, rawData)
	if err != nil {
		return fmt.Errorf("failed to save weather data: %w", err)
	}

	return nil
}

func (s *WeatherStorageDB) GetCitiesWithWeather(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT city_name FROM weather ORDER BY city_name`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get cities with weather: %w", err)
	}
	defer rows.Close()

	var cities []string
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err != nil {
			return nil, fmt.Errorf("failed to scan city: %w", err)
		}
		cities = append(cities, city)
	}

	return cities, nil
}

func (s *WeatherStorageDB) GetCityForecast(ctx context.Context, city string) (*model.CityForecast, error) {
	query := `SELECT date, temp, data FROM weather WHERE city_name = $1 ORDER BY date`
	rows, err := s.db.Query(ctx, query, city)
	if err != nil {
		return nil, fmt.Errorf("failed to get forecast for city %s: %w", city, err)
	}
	defer rows.Close()

	var forecasts []model.Weather
	for rows.Next() {
		var forecast model.Weather
		var rawData []byte
		err := rows.Scan(&forecast.DateTime, &forecast.Temperature, &rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to scan forecast for city %s: %w", city, err)
		}
		err = json.Unmarshal(rawData, &forecast)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal forecast data for city %s: %w", city, err)
		}
		forecasts = append(forecasts, forecast)
	}

	if len(forecasts) == 0 {
		return nil, fmt.Errorf("no forecasts found for city %s", city)
	}

	return &model.CityForecast{
		City:      city,
		Forecasts: forecasts,
	}, nil
}

func (s *WeatherStorageDB) GetWeatherByDateTime(ctx context.Context, city string, dateTime string) (*model.WeatherData, error) {
	query := `SELECT temp, data FROM weather WHERE city_name = $1 AND date = $2`
	s.logger.Infof("Executing query: %s with city: %s and dateTime: %s", query, city, dateTime)
	row := s.db.QueryRow(ctx, query, city, dateTime)

	var temperature float64
	var rawData []byte
	err := row.Scan(&temperature, &rawData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no weather data found for city %s at datetime %s", city, dateTime)
		}
		return nil, fmt.Errorf("failed to get weather data for city %s at datetime %s: %w", city, dateTime, err)
	}

	var weatherDetails model.WeatherDetails
	err = json.Unmarshal(rawData, &weatherDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal weather data for city %s at datetime %s: %w", city, dateTime, err)
	}

	// Формирование структуры WeatherData
	weatherData := &model.WeatherData{
		CityName: city,
		Temp:     temperature,
		Date:     weatherDetails.DtTxt,
		Data:     weatherDetails,
	}

	return weatherData, nil
}
