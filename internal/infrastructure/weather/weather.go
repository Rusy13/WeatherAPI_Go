package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WeatherDetail struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Forecast struct {
	Dt   int64 `json:"dt"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []WeatherDetail `json:"weather"`
	Clouds  struct {
		All int `json:"all"`
	} `json:"clouds"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Visibility int     `json:"visibility"`
	Pop        float64 `json:"pop"`
	Rain       struct {
		H3 float64 `json:"3h"`
	} `json:"rain"`
	Sys struct {
		Pod string `json:"pod"`
	} `json:"sys"`
	DtTxt string `json:"dt_txt"`
}

type WeatherResponse struct {
	Cod     string     `json:"cod"`
	Message float64    `json:"message"`
	Cnt     int        `json:"cnt"`
	List    []Forecast `json:"list"`
	City    struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Coord      Coord  `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Timezone   int    `json:"timezone"`
		Sunrise    int    `json:"sunrise"`
		Sunset     int    `json:"sunset"`
	} `json:"city"`
}

type Coord struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func GetWeatherForecast(lat, lon, apiKey string) (*WeatherResponse, error) {
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%s&lon=%s&appid=%s&units=metric", lat, lon, apiKey)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	var response WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &response, nil
}
