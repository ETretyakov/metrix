// Module "controllers" aggregates all the controllers structures and methods of buisiness logic
// for the service.
package controllers

import (
	"context"
	"sync"

	"metrix/internal/repository"
)

// HealthControllerImpl - the implementation structure for the HealthController that manages
// access to the repositories and mutexes.
type HealthControllerImpl struct {
	repoGroup *repository.Group

	readinessMu sync.RWMutex
	readiness   bool
	livenessMu  sync.RWMutex
	liveness    bool
}

// NewHealthController - the builder function for the HealthControllerImpl.
func NewHealthController(repoGroup *repository.Group) *HealthControllerImpl {
	return &HealthControllerImpl{repoGroup: repoGroup}
}

// SetReadiness - the method that sets "rediness" status for the service.
func (h *HealthControllerImpl) SetReadiness(state bool) {
	h.readinessMu.Lock()
	defer h.readinessMu.Unlock()
	h.readiness = state
}

// SetLiveness - the method that sets "liveness" status for the service.
func (h *HealthControllerImpl) SetLiveness(state bool) {
	h.livenessMu.Lock()
	defer h.livenessMu.Unlock()
	h.liveness = state
}

// ReadinessState - the method that gets "readiness" status for the service.
func (h *HealthControllerImpl) ReadinessState() bool {
	h.readinessMu.RLock()
	defer h.readinessMu.RUnlock()
	return h.readiness
}

// LivenessState - the method that gets "liveness" status for the service.
func (h *HealthControllerImpl) LivenessState() bool {
	h.livenessMu.RLock()
	defer h.livenessMu.RUnlock()
	return h.liveness
}

// PingDB - the method that checks database access in the runtime.
func (h *HealthControllerImpl) PingDB() bool {
	return h.repoGroup.MetricRepo.PingDB(context.Background())
}
