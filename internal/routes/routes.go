package routes

import (
	"net/http"

	"WbTest/internal/middleware"
	"WbTest/internal/weather/delivery"
	"github.com/gorilla/mux"
)

func GetRouter(weatherHandlers *delivery.WeatherDelivery, userHandlers *delivery.UserHandler, mw *middleware.Middleware) *mux.Router {
	router := mux.NewRouter()
	assignRoutes(router, weatherHandlers, userHandlers)
	assignMiddleware(router, mw)
	return router
}

func assignRoutes(router *mux.Router, weatherHandlers *delivery.WeatherDelivery, userHandlers *delivery.UserHandler) {
	router.HandleFunc("/cities", weatherHandlers.GetCities).Methods(http.MethodGet)
	router.HandleFunc("/city/{city}/forecast", weatherHandlers.GetCityForecast).Methods(http.MethodGet)
	router.HandleFunc("/city/{city}/weather/{datetime}", weatherHandlers.GetWeatherByDateTime).Methods(http.MethodGet)
	router.HandleFunc("/register", userHandlers.RegisterUser).Methods(http.MethodPost)
	router.HandleFunc("/login", userHandlers.LoginUser).Methods(http.MethodPost)
	router.HandleFunc("/favorite", userHandlers.AddFavoriteCity).Methods(http.MethodPost)
}

func assignMiddleware(router *mux.Router, mw *middleware.Middleware) {
	router.Use(mw.AccessLog)
	//router.Use(mw.Auth)
}

func serveHTMLFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../templates/index.html")
}
