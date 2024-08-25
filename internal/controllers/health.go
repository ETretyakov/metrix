package controllers

import (
	"context"
	"sync"

	"metrix/internal/repository"
)

type HealthControllerImpl struct {
	repoGroup *repository.Group

	readinessMu sync.RWMutex
	readiness   bool
	livenessMu  sync.RWMutex
	liveness    bool
}

func NewHealthController(repoGroup *repository.Group) *HealthControllerImpl {
	return &HealthControllerImpl{repoGroup: repoGroup}
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

func (h *HealthControllerImpl) PingDB() bool {
	return h.repoGroup.MetricRepo.PingDB(context.Background())
}
