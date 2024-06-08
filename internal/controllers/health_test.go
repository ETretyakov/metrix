package controllers

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthControllerLivenessState(t *testing.T) {
	controller := NewHealthController()

	assert.Equal(t, false, controller.LivenessState())

	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			controller.SetLiveness(true)
		}()
	}
	wg.Wait()

	livenessState := controller.LivenessState()
	assert.Equal(t, true, livenessState)
}

func TestHealthControllerReadinessState(t *testing.T) {
	controller := NewHealthController()

	assert.Equal(t, false, controller.ReadinessState())

	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			controller.SetReadiness(true)
		}()
	}
	wg.Wait()

	livenessState := controller.ReadinessState()
	assert.Equal(t, true, livenessState)
}
