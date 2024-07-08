package delivery

import (
	"WbTest/internal/order/model"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (d *WeatherDelivery) GetCityForecast(w http.ResponseWriter, r *http.Request) {
	city := mux.Vars(r)["city"]

	forecast, err := d.service.GetCityForecast(r.Context(), city)
	if err != nil {
		d.logger.Errorf("error getting forecast for city %s: %v", city, err)
		http.Error(w, fmt.Sprintf("error getting forecast for city %s", city), http.StatusInternalServerError)
		return
	}

	// Get country for the city
	country, err := d.storage.GetCountryForCity(r.Context(), city)
	if err != nil {
		d.logger.Errorf("error getting country for city %s: %v", city, err)
		http.Error(w, fmt.Sprintf("error getting country for city %s", city), http.StatusInternalServerError)
		return
	}

	type ForecastResponse struct {
		Country     string          `json:"country"`
		City        string          `json:"city"`
		AverageTemp float64         `json:"average_temperature"`
		Forecasts   []model.Weather `json:"forecasts"`
	}

	// Calculate average temperature
	var totalTemp float64
	var validForecasts []model.Weather
	for _, f := range forecast.Forecasts {
		totalTemp += f.Temperature
		validForecasts = append(validForecasts, model.Weather{
			DateTime:    f.DateTime,
			Temperature: f.Temperature,
		})
	}
	averageTemp := totalTemp / float64(len(validForecasts))

	response := ForecastResponse{
		Country:     country,
		City:        forecast.City,
		AverageTemp: averageTemp,
		Forecasts:   validForecasts,
	}

	if response.Country == "" {
		response.Country = "Unknown"
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		d.logger.Errorf("error encoding JSON response: %v", err)
		http.Error(w, "error encoding JSON response", http.StatusInternalServerError)
		return
	}
}
