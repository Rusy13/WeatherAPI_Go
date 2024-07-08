package delivery

import (
	"net/http"

	"WbTest/internal/pkg/response"
)

func (d *WeatherDelivery) GetCities(w http.ResponseWriter, r *http.Request) {
	cities, err := d.service.GetCitiesWithWeather(r.Context())
	if err != nil {
		d.logger.Errorf("error getting cities: %v", err)
		response.WriteResponse(w, response.Error{Err: response.ErrInternal.Error()}, http.StatusInternalServerError, d.logger)
		return
	}

	response.WriteResponse(w, cities, http.StatusOK, d.logger)
}
