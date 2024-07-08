package storage

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//
//	"WbTest/internal/infrastructure/database/postgres/database"
//	"WbTest/internal/order/model"
//)
//
//// WeatherStorageImpl реализует интерфейс WeatherStorage.
//type WeatherStorageImpl struct {
//	db *database.PGDatabase
//}
//
//// NewWeatherStorage создает новый экземпляр WeatherStorageImpl.
//func NewWeatherStorage(db *database.PGDatabase) *WeatherStorageImpl {
//	return &WeatherStorageImpl{db: db}
//}
//
//// SaveWeather сохраняет данные о погоде в базу данных.
//func (s *WeatherStorageImpl) SaveWeather(ctx context.Context, city string, weather model.Weather) error {
//	query := `
//		INSERT INTO weather (city_name, date, temperature, humidity, wind_speed, weather_description, raw_data)
//		VALUES ($1, $2, $3, $4, $5, $6, $7)
//		ON CONFLICT (city_name, date) DO UPDATE
//		SET temperature = EXCLUDED.temperature,
//		    humidity = EXCLUDED.humidity,
//		    wind_speed = EXCLUDED.wind_speed,
//		    weather_description = EXCLUDED.weather_description,
//		    raw_data = EXCLUDED.raw_data;
//	`
//
//	rawData, err := json.Marshal(weather)
//	if err != nil {
//		return fmt.Errorf("failed to marshal weather data: %w", err)
//	}
//
//	_, err = s.db.Exec(ctx, query, city, weather.DateTime, weather.Temperature, weather.Humidity, weather.WindSpeed, weather.WeatherDescription, rawData)
//	if err != nil {
//		return fmt.Errorf("failed to save weather data: %w", err)
//	}
//
//	return nil
//}
//
//// GetCitiesWithWeather возвращает список городов, для которых есть данные о погоде.
//func (s *WeatherStorageImpl) GetCitiesWithWeather(ctx context.Context) ([]string, error) {
//	query := `SELECT DISTINCT city_name FROM weather ORDER BY city_name`
//	rows, err := s.db.Query(ctx, query)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get cities with weather: %w", err)
//	}
//	defer rows.Close()
//
//	var cities []string
//	for rows.Next() {
//		var city string
//		if err := rows.Scan(&city); err != nil {
//			return nil, fmt.Errorf("failed to scan city: %w", err)
//		}
//		cities = append(cities, city)
//	}
//
//	return cities, nil
//}
//
//// GetCityForecast возвращает прогноз погоды для указанного города.
//func (s *WeatherStorageDB) GetCityForecast(ctx context.Context, city string) (*model.CityForecast, error) {
//	query := `SELECT date, temperature FROM weather WHERE city_name = $1 ORDER BY date`
//	rows, err := s.db.Query(ctx, query, city)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get forecast for city %s: %w", city, err)
//	}
//	defer rows.Close()
//
//	var forecasts []model.Weather
//	for rows.Next() {
//		var forecast model.Weather
//		if err := rows.Scan(&forecast.DateTime, &forecast.Temperature); err != nil {
//			return nil, fmt.Errorf("failed to scan forecast for city %s: %w", city, err)
//		}
//		forecasts = append(forecasts, forecast)
//	}
//
//	if len(forecasts) == 0 {
//		return nil, fmt.Errorf("no forecasts found for city %s", city)
//	}
//
//	return &model.CityForecast{
//		City:      city,
//		Forecasts: forecasts,
//	}, nil
//}
//
//// GetWeatherByDateTime возвращает данные о погоде для указанного города и времени.
//func (s *WeatherStorageImpl) GetWeatherByDateTime(ctx context.Context, city string, dateTime string) (*model.Weather, error) {
//	query := `SELECT temperature, humidity, wind_speed, weather_description, raw_data FROM weather WHERE city_name = $1 AND date = $2`
//	row := s.db.QueryRow(ctx, query, city, dateTime)
//
//	var weather model.Weather
//	var rawData []byte
//	err := row.Scan(&weather.Temperature, &weather.Humidity, &weather.WindSpeed, &weather.WeatherDescription, &rawData)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get weather data for city %s and datetime %s: %w", city, dateTime, err)
//	}
//
//	err = json.Unmarshal(rawData, &weather)
//	if err != nil {
//		return nil, fmt.Errorf("failed to unmarshal weather data for city %s and datetime %s: %w", city, dateTime, err)
//	}
//
//	return &weather, nil
//}
