package delivery

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"WbTest/internal/order/delivery/dto"
	"WbTest/internal/pkg/response"
)

func (d *OrderDelivery) AddOrder(w http.ResponseWriter, r *http.Request) {
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		d.logger.Errorf("error reading request body: %v", err)
		response.WriteResponse(w, response.Error{Err: response.ErrInternal.Error()}, http.StatusInternalServerError, d.logger)
		return
	}
	defer r.Body.Close()

	var orderDTO dto.AddOrderDTO
	err = json.Unmarshal(rBody, &orderDTO)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			d.logger.Errorf("invalid json: %s", string(rBody))
			response.WriteResponse(w, response.Error{Err: response.ErrInvalidJSON.Error()}, http.StatusBadRequest, d.logger)
			return
		}
		d.logger.Errorf("error unmarshalling request body: %v", err)
		response.WriteResponse(w, response.Error{Err: response.ErrInternal.Error()}, http.StatusInternalServerError, d.logger)
		return
	}

	err = orderDTO.Validate()
	if err != nil {
		d.logger.Errorf("validation errors in adding order: %v", err)
		response.WriteResponse(w, response.Error{Err: err.Error()}, http.StatusBadRequest, d.logger)
		return
	}

	orderToAdd := dto.ConvertToOrder(orderDTO)
	addedOrder, err := d.service.AddOrder(r.Context(), orderToAdd)
	if err != nil {
		if errors.Is(err, ErrDuplicateOrder) {
			d.logger.Errorf("order with ID already exists: %v", orderDTO.OrderUID)
			response.WriteResponse(w, response.Error{Err: err.Error()}, http.StatusBadRequest, d.logger)
			return
		}
		d.logger.Errorf("internal server error in adding order: %v", err)
		response.WriteResponse(w, response.Error{Err: response.ErrInternal.Error()}, http.StatusInternalServerError, d.logger)
		return
	}

	response.WriteResponse(w, dto.OrderResponse{UID: addedOrder.OrderUID}, http.StatusCreated, d.logger)
}
