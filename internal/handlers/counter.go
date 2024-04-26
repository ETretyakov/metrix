package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ETretyakov/metrix/internal/db"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func CounterWidgetUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	log.Info().Msg(fmt.Sprintf("[counter] received variables: %+v", vars))

	name := vars["name"]
	value := vars["value"]

	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		log.Warn().Msg(fmt.Sprintf("[counter] failed to parse value: %s", value))
		w.WriteHeader(http.StatusConflict)
	}

	prevValue := db.MemStorage.Get("counter", name)

	db.MemStorage.Set("counter", name, val+prevValue)

	w.WriteHeader(http.StatusOK)
}
