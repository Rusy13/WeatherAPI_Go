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

type WeatherResponse struct {
	DateTime    string  `json:"date_time"`
	Temperature float64 `json:"temperature"`
}

// WeatherData структура для данных о погоде
type WeatherData struct {
	CityName string         `json:"city_name"`
	Temp     float64        `json:"temp"`
	Date     string         `json:"date"`
	Data     WeatherDetails `json:"data"`
}

// WeatherDetails структура для детальных данных о погоде
type WeatherDetails struct {
	Dt  int64 `json:"dt"`
	Pop int   `json:"pop"`
	Sys struct {
		Pod string `json:"pod"`
	} `json:"sys"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Rain struct {
		ThreeH float64 `json:"3h"`
	} `json:"rain"`
	Wind struct {
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	DtTxt   string `json:"dt_txt"`
	Weather []struct {
		Id          int    `json:"id"`
		Icon        string `json:"icon"`
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Visibility int `json:"visibility"`
}
