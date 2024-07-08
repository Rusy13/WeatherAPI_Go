package routes

import (
	"net/http"

	"WbTest/internal/middleware"
	"WbTest/internal/order/delivery"
	"github.com/gorilla/mux"
)

func GetRouter(handlers *delivery.WeatherDelivery, mw *middleware.Middleware) *mux.Router {
	router := mux.NewRouter()
	assignRoutes(router, handlers)
	assignMiddleware(router, mw)
	return router
}

func assignRoutes(router *mux.Router, handlers *delivery.WeatherDelivery) {
	router.HandleFunc("/cities", handlers.GetCities).Methods(http.MethodGet)
	router.HandleFunc("/city/{city}/forecast", handlers.GetCityForecast).Methods(http.MethodGet)
	router.HandleFunc("/city/{city}/weather/{datetime}", handlers.GetWeatherByDateTime).Methods(http.MethodGet)
}

func assignMiddleware(router *mux.Router, mw *middleware.Middleware) {
	router.Use(mw.AccessLog)
	//router.Use(mw.Auth)
}

func serveHTMLFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../templates/index.html")
}
