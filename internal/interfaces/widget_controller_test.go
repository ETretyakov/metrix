package interfaces_test

import (
	"io"
	"metrix/internal/infrastructure"
	"metrix/internal/interfaces"
	"metrix/internal/logger"
	"metrix/internal/usecases"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWidgetController_Show(t *testing.T) {
	// Dependencies

	storageHandler, err := infrastructure.NewStorageHandler()
	if err != nil {
		logger.Log.Error(err, "testing", "true")
	}

	widgetInteractor := usecases.WidgetInteractor{
		WidgetRepository: &interfaces.WidgetRepository{
			StorageHandler: storageHandler,
		},
	}

	// Setup
	storageHandler.Set("default:gauge:MyGauge", 0)

	// Declare tests
	type fields struct {
		WidgetInteractor usecases.WidgetInteractor
	}
	type args struct {
		w        httptest.ResponseRecorder
		r        *http.Request
		pathVars map[string]string
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Test №1 - Success",
			fields: fields{
				WidgetInteractor: widgetInteractor,
			},
			args: args{
				w: *httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/value/gauge/MyGauge",
					http.NoBody,
				),
				pathVars: map[string]string{
					"widgetType": "gauge",
					"name":       "MyGauge",
				},
			},
			want: want{
				code:        200,
				response:    `0`,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test №2 - Not found",
			fields: fields{
				WidgetInteractor: widgetInteractor,
			},
			args: args{
				w: *httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/value/gauge/MyGaugeNoExists",
					http.NoBody,
				),
				pathVars: map[string]string{
					"widgetType": "gauge",
					"name":       "MyGaugeNoExists",
				},
			},
			want: want{
				code:        404,
				response:    ``,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test №3 - Bad widget type",
			fields: fields{
				WidgetInteractor: widgetInteractor,
			},
			args: args{
				w: *httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/value/gaugeNotExists/MyGaugeNoExists",
					http.NoBody,
				),
				pathVars: map[string]string{
					"widgetType": "gaugeNotExists",
					"name":       "MyGaugeNoExists",
				},
			},
			want: want{
				code:        400,
				response:    ``,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &interfaces.WidgetController{
				WidgetInteractor: tt.fields.WidgetInteractor,
			}
			req := mux.SetURLVars(tt.args.r, tt.args.pathVars)
			wc.Show(&tt.args.w, req)

			res := tt.args.w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			require.NoError(t, err)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
		})
	}
}

func TestWidgetController_Update(t *testing.T) {
	// Dependencies
	storageHandler, err := infrastructure.NewStorageHandler()
	if err != nil {
		logger.Log.Error(err, "testing", "true")
	}

	widgetInteractor := usecases.WidgetInteractor{
		WidgetRepository: &interfaces.WidgetRepository{
			StorageHandler: storageHandler,
		},
	}

	// Setup
	storageHandler.Set("default:gauge:MyGauge", 0)
	storageHandler.Set("default:counter:MyCounter", 0)

	// Declare tests
	type fields struct {
		WidgetInteractor usecases.WidgetInteractor
	}
	type args struct {
		w        httptest.ResponseRecorder
		r        *http.Request
		pathVars map[string]string
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Test №1 - Success (Gauge)",
			fields: fields{
				WidgetInteractor: widgetInteractor,
			},
			args: args{
				w: *httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/value/gauge/MyGauge/100",
					http.NoBody,
				),
				pathVars: map[string]string{
					"widgetType": "gauge",
					"name":       "MyGauge",
					"value":      "100",
				},
			},
			want: want{
				code:        200,
				response:    `100`,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test №2 - Success (Counter)",
			fields: fields{
				WidgetInteractor: widgetInteractor,
			},
			args: args{
				w: *httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/value/counter/MyCounter/100",
					http.NoBody,
				),
				pathVars: map[string]string{
					"widgetType": "counter",
					"name":       "MyCounter",
					"value":      "100",
				},
			},
			want: want{
				code:        200,
				response:    `100`,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test №3 - Success (Counter)",
			fields: fields{
				WidgetInteractor: widgetInteractor,
			},
			args: args{
				w: *httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/value/counter/MyCounter/100",
					http.NoBody,
				),
				pathVars: map[string]string{
					"widgetType": "counter",
					"name":       "MyCounter",
					"value":      "100",
				},
			},
			want: want{
				code:        200,
				response:    `200`,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test №4 - Success (Not Exists)",
			fields: fields{
				WidgetInteractor: widgetInteractor,
			},
			args: args{
				w: *httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/value/counter/MyCounter1/100",
					http.NoBody,
				),
				pathVars: map[string]string{
					"widgetType": "counter",
					"name":       "MyCounter1",
					"value":      "100",
				},
			},
			want: want{
				code:        200,
				response:    `100`,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &interfaces.WidgetController{
				WidgetInteractor: tt.fields.WidgetInteractor,
			}
			req := mux.SetURLVars(tt.args.r, tt.args.pathVars)
			wc.Update(&tt.args.w, req)

			res := tt.args.w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			require.NoError(t, err)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
		})
	}
}

func TestWidgetController_Keys(t *testing.T) {
	// Dependencies
	storageHandler, err := infrastructure.NewStorageHandler()
	if err != nil {
		logger.Log.Error(err, "testing", "true")
	}

	widgetInteractor := usecases.WidgetInteractor{
		WidgetRepository: &interfaces.WidgetRepository{
			StorageHandler: storageHandler,
		},
	}

	// Setup
	storageHandler.Set("default:gauge:MyGauge", 0)

	// Declare tests
	type fields struct {
		WidgetInteractor usecases.WidgetInteractor
	}
	type args struct {
		w        httptest.ResponseRecorder
		r        *http.Request
		pathVars map[string]string
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "Test №1 - Success",
			fields: fields{
				WidgetInteractor: widgetInteractor,
			},
			args: args{
				w: *httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/",
					http.NoBody,
				),
				pathVars: map[string]string{},
			},
			want: want{
				code:        200,
				response:    `{"count":1,"items":[{"namespace":"default","name":"MyGauge","type":"gauge"}]}`,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wc := &interfaces.WidgetController{
				WidgetInteractor: tt.fields.WidgetInteractor,
			}
			req := mux.SetURLVars(tt.args.r, tt.args.pathVars)
			wc.Keys(&tt.args.w, req)

			res := tt.args.w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			require.NoError(t, err)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			if len(tt.want.response) > 0 {
				assert.JSONEq(t, tt.want.response, string(resBody))
			}
		})
	}
}
