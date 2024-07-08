package delivery

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func (d *WeatherDelivery) GetWeatherByDateTime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	city := vars["city"]
	dateTime := vars["datetime"]

	// Проверка формата даты и времени
	parsedTime, err := time.Parse("2006-01-02T15:04:05", dateTime)
	if err != nil {
		d.logger.Errorf("error parsing datetime %s: %v", dateTime, err)
		http.Error(w, "Invalid datetime format", http.StatusBadRequest)
		return
	}

	// Преобразование времени в строку для запроса
	dateTimeStr := parsedTime.Format("2006-01-02 15:04:05")
	d.logger.Infof("Fetching weather data for city %s at datetime %s", city, dateTimeStr)

	weather, err := d.storage.GetWeatherByDateTime(r.Context(), city, dateTimeStr)
	if err != nil {
		d.logger.Errorf("error getting weather data for city %s at datetime %s: %v", city, dateTimeStr, err)
		http.Error(w, "Weather data not found", http.StatusNotFound)
		return
	}

	response, err := json.Marshal(weather)
	if err != nil {
		d.logger.Errorf("error marshaling weather data for city %s at datetime %s: %v", city, dateTimeStr, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
