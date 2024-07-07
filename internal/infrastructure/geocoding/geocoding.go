package geocoding

import (
	dbcity "WbTest/internal/order/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type CityInfo struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

func GetCityInfo(cityName, apiKey string) (dbcity.City, error) {
	var cityInfo dbcity.City

	baseURL := "http://api.openweathermap.org/geo/1.0/direct"
	query := fmt.Sprintf("%s?q=%s&limit=1&appid=%s", baseURL, url.QueryEscape(cityName), apiKey)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(query)
	if err != nil {
		return cityInfo, fmt.Errorf("failed to get city info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return cityInfo, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	var cityResults []struct {
		Name    string  `json:"name"`
		Country string  `json:"country"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
	}
	err = json.NewDecoder(resp.Body).Decode(&cityResults)
	if err != nil {
		return cityInfo, fmt.Errorf("failed to decode city info: %w", err)
	}

	if len(cityResults) == 0 {
		return cityInfo, fmt.Errorf("no results found for city: %s", cityName)
	}

	cityResult := cityResults[0]
	cityInfo = dbcity.City{
		Name:      cityResult.Name,
		Country:   cityResult.Country,
		Latitude:  cityResult.Lat,
		Longitude: cityResult.Lon,
	}

	return cityInfo, nil
}
