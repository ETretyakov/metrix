// Module "handlers" aggregates all the handlers structures and methods for the service.
package handlers

import (
	"net/http"

	"metrix/internal/controllers"
	"metrix/internal/repository"
)

// HealthHandlers - the implementation structure for the HealthHandlers that manages
// access to the related controller.
type HealthHandlers struct {
	controller controllers.HealthController
}

// NewHealthHandlers - the builder function for the HealthHandlers.
func NewHealthHandlers(repoGroup *repository.Group) *HealthHandlers {
	controller := controllers.NewHealthController(repoGroup)
	return &HealthHandlers{
		controller: controller,
	}
}

// SetReadiness - the method that sets "rediness" status for the service via handler.
func (h *HealthHandlers) SetReadiness(state bool) {
	h.controller.SetReadiness(state)
}

// SetLiveness - the method that sets "liveness" status for the service via handler.
func (h *HealthHandlers) SetLiveness(state bool) {
	h.controller.SetLiveness(state)
}

// ReadinessState - the method that gets "readiness" status for the service via handler.
// @Tags Info
// @Summary Query to retrieve service readiness state
// @ID infoReadiness
// @Produce plain
// @Success 200
// @Failure 500 {string} string "Internal Server Error"
// @Router /readiness [get]
func (h *HealthHandlers) ReadinessState(w http.ResponseWriter, r *http.Request) {
	if h.controller.ReadinessState() {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}

// LivenessState - the method that gets "liveness" status for the service via handler.
// @Tags Info
// @Summary Query to retrieve service liveness state
// @ID infoLiveness
// @Produce plain
// @Success 200
// @Failure 500 {string} string "Internal Server Error"
// @Router /liveness [get]
func (h *HealthHandlers) LivenessState(w http.ResponseWriter, r *http.Request) {
	if h.controller.LivenessState() {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}

// PingDB - the method that checks database access in the runtime via handler.
// @Tags Info
// @Summary Query to retrieve service database connection state
// @ID infoPingDB
// @Produce plain
// @Success 200
// @Failure 500 {string} string "Internal Server Error"
// @Router /liveness [get]
func (h *HealthHandlers) PingDB(w http.ResponseWriter, r *http.Request) {
	if h.controller.PingDB() {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}
