package delivery

import (
	"WbTest/internal/pkg/response"
	"github.com/gorilla/mux"
	"net/http"
)

func (d *WeatherDelivery) GetCityForecast(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	city := params["city"]

	forecast, err := d.service.GetCityForecast(r.Context(), city)
	if err != nil {
		d.logger.Errorf("error getting forecast for city %s: %v", city, err)
		response.WriteResponse(w, response.Error{Err: "forecast not found"}, http.StatusNotFound, d.logger)
		return
	}

	response.WriteResponse(w, forecast, http.StatusOK, d.logger)
}
