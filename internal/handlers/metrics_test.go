package handlers

import (
	"bytes"
	"context"
	"metrix/internal/controllers"
	"metrix/internal/repository"
	"metrix/internal/validators"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestMetricsHandlers_Set(t *testing.T) {
	ctx := context.Background()
	repoGroup := repository.NewGroup(ctx, nil, "", 0, false)
	controller := controllers.NewMetricController(repoGroup)
	validator := validators.NewMetricsValidator()

	type fields struct {
		controller controllers.MetricsController
		validator  validators.MetricsValidator
	}
	type args struct {
		w              http.ResponseWriter
		r              *http.Request
		metricType     string
		metricID       string
		metricValue    string
		wantStatusCode int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test 1: Set metric value",
			fields: fields{
				controller: controller,
				validator:  validator,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/update/counter/metric_1/10",
					http.NoBody,
				),
				metricType:     "counter",
				metricID:       "metric_1",
				metricValue:    "10",
				wantStatusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &MetricsHandlers{
				controller: tt.fields.controller,
				validator:  tt.fields.validator,
			}
			tt.args.r = mux.SetURLVars(tt.args.r, map[string]string{
				"type":  tt.args.metricType,
				"id":    tt.args.metricID,
				"value": tt.args.metricValue,
			})
			ctx := tt.args.r.Context()
			h.Set(tt.args.w, tt.args.r.WithContext(ctx))
			if r, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if r.Code != tt.args.wantStatusCode {
					t.Errorf(
						"status codes are different: got=%d want=%d",
						r.Code,
						tt.args.wantStatusCode,
					)
				}
			} else {
				t.Errorf("got different from *httptest.ResponseRecorder struct: %+v", r)
			}
		})
	}
}

func TestMetricsHandlers_SetWithModel(t *testing.T) {
	ctx := context.Background()
	repoGroup := repository.NewGroup(ctx, nil, "", 0, false)
	controller := controllers.NewMetricController(repoGroup)
	validator := validators.NewMetricsValidator()

	type fields struct {
		controller controllers.MetricsController
		validator  validators.MetricsValidator
	}
	type args struct {
		w              http.ResponseWriter
		r              *http.Request
		wantStatusCode int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test 1: Set metric",
			fields: fields{
				controller: controller,
				validator:  validator,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/update/",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`{"id": "metric_1", "type": "counter", "delta": 10}`)
						}(),
					),
				),
				wantStatusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &MetricsHandlers{
				controller: tt.fields.controller,
				validator:  tt.fields.validator,
			}
			h.SetWithModel(tt.args.w, tt.args.r)
			if r, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if r.Code != tt.args.wantStatusCode {
					t.Errorf(
						"status codes are different: got=%d want=%d",
						r.Code,
						tt.args.wantStatusCode,
					)
				}
			} else {
				t.Errorf("got different from *httptest.ResponseRecorder struct: %+v", r)
			}
		})
	}
}

func TestMetricsHandlers_SetMany(t *testing.T) {
	ctx := context.Background()
	repoGroup := repository.NewGroup(ctx, nil, "", 0, false)
	controller := controllers.NewMetricController(repoGroup)
	validator := validators.NewMetricsValidator()

	type fields struct {
		controller controllers.MetricsController
		validator  validators.MetricsValidator
	}
	type args struct {
		w              http.ResponseWriter
		r              *http.Request
		wantStatusCode int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test 1: Set metrics",
			fields: fields{
				controller: controller,
				validator:  validator,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/updates/",
					bytes.NewBuffer(
						func() []byte {
							return []byte(`[{"id": "metric_1", "type": "counter", "delta": 10}]`)
						}(),
					),
				),
				wantStatusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &MetricsHandlers{
				controller: tt.fields.controller,
				validator:  tt.fields.validator,
			}
			h.SetMany(tt.args.w, tt.args.r)
			if r, ok := tt.args.w.(*httptest.ResponseRecorder); ok {
				if r.Code != tt.args.wantStatusCode {
					t.Errorf(
						"status codes are different: got=%d want=%d",
						r.Code,
						tt.args.wantStatusCode,
					)
				}
			} else {
				t.Errorf("got different from *httptest.ResponseRecorder struct: %+v", r)
			}
		})
	}
}
