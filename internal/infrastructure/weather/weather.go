package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Forecast struct {
	Dt   int64 `json:"dt"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	List []Forecast `json:"list"`
}

func GetWeatherForecast(lat, lon, apiKey string) (Forecast, error) {
	var forecast Forecast

	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%s&lon=%s&appid=%s&units=metric", lat, lon, apiKey)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return forecast, fmt.Errorf("failed to get weather forecast: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return forecast, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&forecast)
	if err != nil {
		return forecast, fmt.Errorf("failed to decode weather forecast: %w", err)
	}

	return forecast, nil
}
