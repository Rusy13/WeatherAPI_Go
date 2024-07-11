package storage

import (
	PG "WbTest/internal/infrastructure/database/postgres/database"
	"context"
	"fmt"
)

type City struct {
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func SaveCity(db *PG.PGDatabase, city City) error {
	query := `
        INSERT INTO cities (name, country, latitude, longitude)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (name) DO UPDATE 
        SET country = EXCLUDED.country,
            latitude = EXCLUDED.latitude,
            longitude = EXCLUDED.longitude;
    `
	_, err := db.Exec(context.Background(), query, city.Name, city.Country, city.Latitude, city.Longitude)
	if err != nil {
		return fmt.Errorf("failed to save city: %w", err)
	}
	return nil
}

func GetCities(ctx context.Context, db *PG.PGDatabase) ([]City, error) {
	rows, err := db.Query(ctx, "SELECT name, country, latitude, longitude FROM cities ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []City
	for rows.Next() {
		var city City
		if err := rows.Scan(&city.Name, &city.Country, &city.Latitude, &city.Longitude); err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}
	return cities, nil
}
