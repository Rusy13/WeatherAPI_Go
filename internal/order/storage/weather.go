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
	CityName string
	Temp     float64
	Date     string
	Data     string
}

func SaveWeather(db *PG.PGDatabase, cityName string, forecast []weather.Forecast) error {
	query := `
        INSERT INTO weather (city_name, temp, date, data)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (city_name, date) DO UPDATE 
        SET temp = EXCLUDED.temp,
            data = EXCLUDED.data;
    `
	for _, f := range forecast {
		date := time.Unix(f.Dt, 0).Format("2006-01-02")
		data, err := json.Marshal(f)
		if err != nil {
			return fmt.Errorf("failed to marshal forecast data: %w", err)
		}

		_, err = db.Exec(context.Background(), query, cityName, f.Main.Temp, date, data)
		if err != nil {
			return fmt.Errorf("failed to save weather forecast: %w", err)
		}
	}
	return nil
}

func GetWeather(db *PG.PGDatabase, cityName, date string) (Weather, error) {
	query := `SELECT city_name, temp, date, data FROM weather WHERE city_name = $1 AND date = $2`
	row := db.QueryRow(context.Background(), query, cityName, date)

	var weather Weather
	if err := row.Scan(&weather.CityName, &weather.Temp, &weather.Date, &weather.Data); err != nil {
		if err == sql.ErrNoRows {
			return weather, nil
		}
		return weather, err
	}

	return weather, nil
}
