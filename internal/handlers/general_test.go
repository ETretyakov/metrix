package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test №1",
			want: want{
				code:        200,
				response:    `{"status": true, "message": "I'm fine, thank you!"}`,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/health",
				nil,
			)

			w := httptest.NewRecorder()
			HealthHandler(w, req)

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
