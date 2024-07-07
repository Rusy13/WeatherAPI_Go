package main

import (
	"WbTest/internal/infrastructure/geocoding"
	"WbTest/internal/infrastructure/weather"
	dbcity "WbTest/internal/order/storage"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	//"os"

	"WbTest/internal/infrastructure/database/postgres/database"
	"WbTest/internal/infrastructure/database/redis"

	"WbTest/internal/middleware"
	"WbTest/internal/order/delivery"
	serviceOrder "WbTest/internal/order/service"
	DbCities "WbTest/internal/order/storage"
	storageOrder "WbTest/internal/order/storage/database"
	"WbTest/internal/routes"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type GeocodingResponse struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
}

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
	//err = godotenv.Load(".env.testing")
	//if err != nil {
	//	logger.Fatalf("error in getting env: %s", err)
	//}
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
	//-----------------------------------------------------------------------------------------------------------------------------------------------------
	cities := []string{"London", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose", "Austin", "Jacksonville", "Fort Worth", "Columbus", "Charlotte", "San Francisco", "Indianapolis", "Seattle", "Denver", "Washington"}

	for _, cityName := range cities {
		//city, err := geocoding.GetCityInfo(cityName, config.LoadConfig().GeocodingAPIKey)
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
	//------------------------------------------------------------------------------------------------------------
	cityList, err := DbCities.GetCities(context.Background(), db)
	if err != nil {
		log.Fatal(err)
	}

	for _, city := range cityList {
		forecast, err := weather.GetWeatherForecast(fmt.Sprintf("%f", city.Latitude), fmt.Sprintf("%f", city.Longitude), "a0685ac1c3e28d3319a052a1d6897687")
		if err != nil {
			log.Printf("Failed to get forecast for city %s: %v", city.Name, err)
			continue
		}

		err = DbCities.SaveWeather(db, city.Name, forecast.List)
		if err != nil {
			log.Printf("Failed to save forecast for city %s: %v", city.Name, err)
		}
	}

	log.Println("Weather data updated successfully")
	//------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	redisConn, err := redis.Init()
	if err != nil {
		logger.Fatalf("error on connection to redis: %v", err)
	}
	defer func() {
		err = redisConn.Close()
		if err != nil {
			logger.Infof("error on redis close: %s", err.Error())
		}
	}()

	stOrder := storageOrder.New(db, redisConn, logger)
	err = stOrder.RestoreCacheFromDB() // Вызов функции восстановления кэша из БД
	if err != nil {
		logger.Fatalf("error restoring cache from DB: %v", err)
	}
	svOrder := serviceOrder.New(stOrder)
	d := delivery.New(svOrder, logger)

	mw := middleware.New(logger)
	router := routes.GetRouter(d, mw)

	//////////////////////////////////////////////////NATS//////////////////////////////////////////////
	//natsConn, err := stan.Connect("test-cluster", "order-service", stan.NatsURL("nats://localhost:4222"))
	//if err != nil {
	//	logger.Fatalf("error connecting to NATS Streaming: %v", err)
	//}
	//defer natsConn.Close()
	//
	//natsHandler := nats.NewNatsHandler(natsConn, svOrder)
	//err = natsHandler.Subscribe("order-channel")
	//if err != nil {
	//	logger.Fatalf("error subscribing to NATS channel: %v", err)
	//}
	//defer natsHandler.Close()
	//////////////////////////////////////////////////NATS//////////////////////////////////////////////
	//port := os.Getenv("APP_PORT")

	port := "8000"
	addr := ":" + port
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	logger.Fatal(http.ListenAndServe(addr, router))
}

func saveCityToDB(db *sql.DB, city GeocodingResponse) {
	insertCitySQL := `INSERT INTO cities (name, country, latitude, longitude) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertCitySQL)
	if err != nil {
		log.Fatalln(err)
	}
	defer statement.Close()

	_, err = statement.Exec(city.Name, city.Country, city.Lat, city.Lon)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Saved city: %s, %s\n", city.Name, city.Country)
}
