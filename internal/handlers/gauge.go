package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ETretyakov/metrix/internal/db"

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
		w.WriteHeader(http.StatusConflict)
	}

	db.MemStorage.Set("gauge", name, val)

	w.WriteHeader(http.StatusOK)
}
