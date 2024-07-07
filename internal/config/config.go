package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	GeocodingAPIKey string
	WeatherAPIKey   string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	config := Config{
		GeocodingAPIKey: os.Getenv("GEOCODING_API_KEY"),
		WeatherAPIKey:   os.Getenv("WEATHER_API_KEY"),
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASS"),
		DBName:          os.Getenv("DB_NAME"),
	}

	return config
}
