package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ETretyakov/metrix/internal/schemas"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	res := schemas.StatusMessage{
		Status: true,
		Msg:    "I'm fine, thank you!",
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}

func UnknownMetricHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}
