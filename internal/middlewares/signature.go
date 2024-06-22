package middlewares

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"metrix/pkg/logger"

	"net/http"
)

var signKey string

func SetSignKey(key string) {
	signKey = key
}

func SignatureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hashSum := r.Header.Get("HashSHA256")
		if r.Method == http.MethodPost && len(hashSum) != 0 {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Warn(
					context.TODO(),
					"failed to read body for checking signature",
					"url", r.URL,
					"method", r.Method,
				)
			}
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			h := hmac.New(sha256.New, []byte(signKey))
			h.Write(bodyBytes)
			signature := hex.EncodeToString(h.Sum(nil))

			w.Header().Add("HashSHA256", signature)

			if signature != hashSum {
				logger.Warn(
					context.TODO(),
					fmt.Sprintf(
						"signKey=%s",
						signKey,
					),
				)
				logger.Warn(
					context.TODO(),
					fmt.Sprintf(
						"received body=%s",
						bodyBytes,
					),
				)
				logger.Warn(
					context.TODO(),
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
