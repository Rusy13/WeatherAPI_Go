package response

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func WriteResponse(w http.ResponseWriter, responseToWrite interface{}, statusCode int, logger *zap.SugaredLogger) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Connection", "keep-alive")
	responseJSON, err := json.Marshal(responseToWrite)
	if err != nil {
		logger.Errorf("error in marshalling response: %v", err)
		WriteResponse(w, []byte(`{"error":"internal error"}`), http.StatusInternalServerError, logger)
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(responseJSON)
	if err != nil {
		logger.Errorf("error in writing response: %v", err)
	}
}
