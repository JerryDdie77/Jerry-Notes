package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("GET /health")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var response struct {
		Code   int    `json:"code"`
		Status string `json:"status"`
	}
	response.Code = http.StatusOK
	response.Status = "OK"
	json.NewEncoder(w).Encode(response)
}
