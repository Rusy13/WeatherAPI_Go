package model

import "time"

// Weather представляет данные о погоде для конкретного времени.
type Weather struct {
	DateTime           time.Time `json:"date_time"`
	Temperature        float64   `json:"temperature"`
	Humidity           float64   `json:"humidity"`
	WindSpeed          float64   `json:"wind_speed"`
	WeatherDescription string    `json:"weather_description"`
}

// CityForecast представляет прогноз погоды для города.
type CityForecast struct {
	City      string    `json:"city"`
	Forecasts []Weather `json:"forecasts"`
}
