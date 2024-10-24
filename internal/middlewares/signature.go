package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"metrix/pkg/logger"

	"net/http"
)

var signKey string = ""

// SetSignKey - the function that sets string for a global variable.
func SetSignKey(key string) {
	signKey = key
}

// SignatureMiddleware - the net/http middleware function to signt http content.
func SignatureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hashSum := r.Header.Get("HashSHA256")
		if r.Method == http.MethodPost && hashSum != "" && signKey != "" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Warn(
					r.Context(),
					"failed to read body for checking signature",
					"url", r.URL,
					"method", r.Method,
				)
			}
			if err := r.Body.Close(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			h := hmac.New(sha256.New, []byte(signKey))
			h.Write(bodyBytes)
			signature := hex.EncodeToString(h.Sum(nil))

			w.Header().Add("HashSHA256", signature)

			if signature != hashSum {
				logger.Warn(
					r.Context(),
					fmt.Sprintf(
						"wrong signature calc=%s got=%s",
						signature,
						hashSum,
					),
				)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
