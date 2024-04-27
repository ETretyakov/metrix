package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGaugeWidgetUpdateHandler(t *testing.T) {
	type want struct {
		code        int
		contentType string
		response    string
	}

	tests := []struct {
		name       string
		widgetName string
		value      string
		want       want
	}{
		{
			name:       "positive test №1",
			widgetName: "TestGauge1",
			value:      "1",
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
				response:    `{"value": 1}`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := fmt.Sprintf("/update/gauge/%s/%s", test.widgetName, test.value)
			req := httptest.NewRequest(
				http.MethodPost,
				url,
				nil,
			)

			vars := map[string]string{
				"name":  test.widgetName,
				"value": test.value,
			}
			req = mux.SetURLVars(req, vars)

			w := httptest.NewRecorder()
			CounterWidgetUpdateHandler(w, req)

			res := w.Result()
			assert.Equal(t, res.StatusCode, test.want.code)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.JSONEq(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}

}
