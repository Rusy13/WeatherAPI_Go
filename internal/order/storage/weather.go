package storage

import (
	PG "WbTest/internal/infrastructure/database/postgres/database"
	"WbTest/internal/infrastructure/weather"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Weather struct {
	CityName    string
	Temperature float64
	DateTime    string
	Data        string
}

func SaveWeatherJson(db *PG.PGDatabase, cityName string, bodyBytes []byte) error {
	var response weather.WeatherResponse
	err := json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal forecast JSON for city %s: %v", cityName, err)
	}

	for _, forecast := range response.List {
		date := time.Unix(forecast.Dt, 0).Format("2006-01-02 15:04:05")
		data, err := json.Marshal(forecast)
		if err != nil {
			return fmt.Errorf("failed to marshal forecast data: %v", err)
		}

		query := `
            INSERT INTO weather (city_name, date, temp, data)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (city_name, date) DO UPDATE 
            SET temp = EXCLUDED.temp,
                data = EXCLUDED.data;
        `

		_, err = db.Exec(context.Background(), query, cityName, date, forecast.Main.Temp, data)
		if err != nil {
			return fmt.Errorf("failed to save weather forecast: %v", err)
		}
	}

	return nil
}

func GetWeather(db *PG.PGDatabase, cityName, date string) (Weather, error) {
	query := `SELECT city_name, temp, date, data FROM weather WHERE city_name = $1 AND date = $2`
	row := db.QueryRow(context.Background(), query, cityName, date)

	var weather Weather
	if err := row.Scan(&weather.CityName, &weather.Temperature, &weather.DateTime, &weather.Data); err != nil {
		if err == sql.ErrNoRows {
			return weather, nil
		}
		return weather, err
	}

	return weather, nil
}
