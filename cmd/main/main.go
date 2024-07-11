package main

import (
	"WbTest/internal/infrastructure/geocoding"
	"WbTest/internal/infrastructure/weather"
	dbcity "WbTest/internal/weather/storage"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"WbTest/internal/infrastructure/database/postgres/database"
	"WbTest/internal/middleware"
	"WbTest/internal/routes"
	"WbTest/internal/weather/delivery"
	serviceOrder "WbTest/internal/weather/service"
	DbCities "WbTest/internal/weather/storage"
	storageOrder "WbTest/internal/weather/storage/database"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
		log.Println("Error is-----------------------", err)
	}

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error in logger initialization: %v", err)
	}
	logger := zapLogger.Sugar()
	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Printf("error in logger sync: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := database.New(ctx)
	if err != nil {
		logger.Fatalf("error in database init: %v", err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			logger.Errorf("error in closing db")
		}
	}()

	cities := []string{"London", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose", "Austin", "Jacksonville", "Fort Worth", "Columbus", "Charlotte", "San Francisco", "Indianapolis", "Seattle", "Denver", "Washington"}

	for _, cityName := range cities {
		city, err := geocoding.GetCityInfo(cityName, "a0685ac1c3e28d3319a052a1d6897687")

		if err != nil {
			log.Printf("Failed to get info for city %s: %v", cityName, err)
			continue
		}

		log.Printf("City info: %v", city)

		if err := dbcity.SaveCity(db, city); err != nil {
			log.Printf("Failed to save city %s: %v", cityName, err)
			continue
		}
	}

	// Фоновый процесс для обновления данных о погоде
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		log.Printf("NewTicker")
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				updateWeatherData(db)
			case <-ctx.Done():
				return
			}
		}
	}()

	stOrder := storageOrder.New(db, logger)
	userStorage := storageOrder.NewUserStorageDB(db, logger)

	svOrder := serviceOrder.New(stOrder)
	userService := serviceOrder.NewUserService(userStorage)

	d := delivery.New(svOrder, stOrder, logger)
	userHandlers := delivery.NewUserHandler(userService)

	mw := middleware.New(logger)
	router := routes.GetRouter(d, userHandlers, mw)

	port := "8000"
	addr := ":" + port
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	logger.Fatal(http.ListenAndServe(addr, router))
}

func updateWeatherData(db *database.PGDatabase) {
	cityList, err := DbCities.GetCities(context.Background(), db)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	apiKey := "a0685ac1c3e28d3319a052a1d6897687"

	//Дополнительно. Распараллелить процесс получения информации о
	//погоде из внешнего API

	for _, city := range cityList {
		wg.Add(1)
		go func(city dbcity.City) {
			defer wg.Done()

			response, err := weather.GetWeatherForecast(fmt.Sprintf("%f", city.Latitude), fmt.Sprintf("%f", city.Longitude), apiKey)
			if err != nil {
				log.Printf("Failed to get forecast for city %s: %v", city.Name, err)
				return
			}

			forecastsBytes, err := json.Marshal(response)
			if err != nil {
				log.Printf("Failed to marshal forecast data for city %s: %v", city.Name, err)
				return
			}

			err = DbCities.SaveWeatherJson(db, city.Name, forecastsBytes)
			if err != nil {
				log.Printf("Failed to save forecast for city %s: %v", city.Name, err)
			}
		}(city)
	}

	wg.Wait()
	log.Println("Weather data updated successfully")
}
