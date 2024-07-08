package model

import (
	"time"
)

type Weather struct {
	DateTime    time.Time `json:"date_time"`
	Temperature float64   `json:"temperature"`
}

// NullFloat64 представляет float64 значение с возможностью быть null
type NullFloat64 struct {
	Valid   bool
	Float64 float64
}

type CityForecast struct {
	City           string    `json:"city"`
	Country        string    `json:"country"`
	AverageTemp    float64   `json:"average_temperature"`
	AvailableDates []string  `json:"available_dates"`
	Forecasts      []Weather `json:"forecasts"`
}
