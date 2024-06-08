package controllers

import "sync"

type HealthControllerImpl struct {
	readinessMu sync.RWMutex
	readiness   bool
	livenessMu  sync.RWMutex
	liveness    bool
}

func NewHealthController() *HealthControllerImpl {
	return &HealthControllerImpl{}
}

func (h *HealthControllerImpl) SetReadiness(state bool) {
	h.readinessMu.Lock()
	defer h.readinessMu.Unlock()
	h.readiness = state
}

func (h *HealthControllerImpl) SetLiveness(state bool) {
	h.livenessMu.Lock()
	defer h.livenessMu.Unlock()
	h.liveness = state
}

func (h *HealthControllerImpl) ReadinessState() bool {
	h.readinessMu.RLock()
	defer h.readinessMu.RUnlock()
	return h.readiness
}

func (h *HealthControllerImpl) LivenessState() bool {
	h.livenessMu.RLock()
	defer h.livenessMu.RUnlock()
	return h.liveness
}
