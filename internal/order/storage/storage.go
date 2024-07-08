package storage

import (
	"WbTest/internal/order/model"
	"context"
)

// WeatherStorage определяет методы для работы с данными о погоде.
type WeatherStorage interface {
	SaveWeather(ctx context.Context, city string, weather model.Weather) error
	GetCitiesWithWeather(ctx context.Context) ([]string, error)
	GetCityForecast(ctx context.Context, city string) (*model.CityForecast, error)
	GetWeatherByDateTime(ctx context.Context, city string, dateTime string) (*model.Weather, error)
}
