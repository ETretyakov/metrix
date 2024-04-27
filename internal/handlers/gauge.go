package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ETretyakov/metrix/internal/db"
	"github.com/ETretyakov/metrix/internal/schemas"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func GaugeWidgetUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	log.Info().Msg(fmt.Sprintf("[gauge] received variables: %+v", vars))

	name := vars["name"]
	value := vars["value"]

	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		log.Warn().Msg(fmt.Sprintf("[gauge] failed to parse value: %s", value))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db.MemStorage.Set("gauge", name, val)

	res := schemas.WidgetResponse{
		Value: db.MemStorage.Get("counter", name),
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}
