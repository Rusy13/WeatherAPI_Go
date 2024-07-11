package service

import (
	"context"
	"time"

	"WbTest/internal/weather/model"
	"WbTest/internal/weather/storage"
)

// WeatherService определяет методы для работы с данными о погоде.
type WeatherService interface {
	GetCitiesWithWeather(ctx context.Context) ([]string, error)
	GetCityForecast(ctx context.Context, city string) (*model.CityForecast, error)
	GetWeatherByDateTime(ctx context.Context, city string, dateTime time.Time) (*model.WeatherData, error)
}

// WeatherServiceImpl реализует интерфейс WeatherService.
type WeatherServiceImpl struct {
	storage storage.WeatherStorage
}

// NewWeatherService создает новый экземпляр WeatherServiceImpl.
func New(storage storage.WeatherStorage) *WeatherServiceImpl {
	return &WeatherServiceImpl{
		storage: storage,
	}
}

// GetCitiesWithWeather возвращает список городов, для которых есть данные о погоде.
func (s *WeatherServiceImpl) GetCitiesWithWeather(ctx context.Context) ([]string, error) {
	return s.storage.GetCitiesWithWeather(ctx)
}

// GetCityForecast возвращает прогноз погоды для указанного города.
func (s *WeatherServiceImpl) GetCityForecast(ctx context.Context, city string) (*model.CityForecast, error) {
	return s.storage.GetCityForecast(ctx, city)
}

// GetWeatherByDateTime возвращает данные о погоде для указанного города и времени.
func (s *WeatherServiceImpl) GetWeatherByDateTime(ctx context.Context, city string, dateTime time.Time) (*model.WeatherData, error) {
	// Приведение времени к формату, который используется в хранилище (если необходимо)
	dateTimeStr := dateTime.Format("2006-01-02 15:04:05")
	return s.storage.GetWeatherByDateTime(ctx, city, dateTimeStr)
}
