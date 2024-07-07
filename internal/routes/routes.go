package routes

import (
	"net/http"

	"WbTest/internal/middleware"
	"WbTest/internal/order/delivery"
	"github.com/gorilla/mux"
)

func GetRouter(handlers *delivery.OrderDelivery, mw *middleware.Middleware) *mux.Router {
	router := mux.NewRouter()
	assignRoutes(router, handlers)
	assignMiddleware(router, mw)
	return router
}

func assignRoutes(router *mux.Router, handlers *delivery.OrderDelivery) {
	router.HandleFunc("/order", handlers.AddOrder).Methods(http.MethodPost)
	router.HandleFunc("/order/{id}", handlers.GetOrder).Methods(http.MethodGet)
	router.HandleFunc("/", serveHTMLFile) // Добавьте этот обработчик для обслуживания HTML файла

}

func assignMiddleware(router *mux.Router, mw *middleware.Middleware) {
	router.Use(mw.AccessLog)
	//router.Use(mw.Auth)
}

func serveHTMLFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../../templates/index.html")
}
