package delivery

import (
	"WbTest/internal/pkg/response"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (d *OrderDelivery) GetOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	orderID := params["id"]
	log.Println("DELIV =========================== ", orderID)

	//filters, err := d.getFilters(params, w)
	//if err != nil {
	//	return
	//}

	order, err := d.service.GetUserOrder(r.Context(), orderID)
	if err != nil {
		d.logger.Errorf("error getting order: %v", err)
		response.WriteResponse(w, response.Error{Err: "order not found"}, http.StatusNotFound, d.logger)
		return
	}

	response.WriteResponse(w, order, http.StatusOK, d.logger)
}
