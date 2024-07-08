package delivery

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"WbTest/internal/pkg/response"
)

func (d *WeatherDelivery) GetWeatherByDateTime(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	city := params["city"]
	dateTimeStr := params["datetime"]

	dateTime, err := time.Parse("2006-01-02T15:04:05", dateTimeStr)
	if err != nil {
		d.logger.Errorf("invalid datetime format: %v", err)
		response.WriteResponse(w, response.Error{Err: "invalid datetime format"}, http.StatusBadRequest, d.logger)
		return
	}

	weather, err := d.service.GetWeatherByDateTime(r.Context(), city, dateTime)
	if err != nil {
		d.logger.Errorf("error getting weather data for city %s at datetime %s: %v", city, dateTimeStr, err)
		response.WriteResponse(w, response.Error{Err: "weather data not found"}, http.StatusNotFound, d.logger)
		return
	}

	response.WriteResponse(w, weather, http.StatusOK, d.logger)
}
