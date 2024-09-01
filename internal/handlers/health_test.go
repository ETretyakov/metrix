// Module "handlers" aggregates all the handlers structures and methods for the service.
package handlers

import (
	"context"
	"metrix/internal/controllers"
	"metrix/internal/repository"
	"testing"
)

func TestHealthHandlers_SetReadiness(t *testing.T) {
	ctx := context.Background()
	repoGroup := repository.NewGroup(ctx, nil, "", 0, false)
	controller := controllers.NewHealthController(repoGroup)

	type fields struct {
		controller controllers.HealthController
	}
	type args struct {
		state bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Test 1: Set Readiness True",
			fields: fields{
				controller: controller,
			},
			args: args{state: true},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HealthHandlers{
				controller: tt.fields.controller,
			}
			h.SetReadiness(tt.args.state)

			state := h.controller.ReadinessState()
			if tt.want != state {
				t.Errorf("HealthHandlers.SetReadiness() = %v, want %v", state, tt.want)
			}
		})
	}
}
