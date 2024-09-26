package middlewares

import (
	"bytes"
	"io"
	"metrix/pkg/crypto"
	"metrix/pkg/logger"
	"net/http"

	"github.com/pkg/errors"
)

var decryption *crypto.Decryption

func InitDecryption(privateKeyPath string) error {
	dcr, err := crypto.NewDecryption(privateKeyPath)
	if err != nil {
		return errors.Wrap(err, "failed to init decription")
	}
	decryption = dcr
	return nil
}

func DecryptionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-encrypted") == "true" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Warn(
					r.Context(),
					"failed to read body",
					"url", r.URL,
					"method", r.Method,
				)
			}
			if err := r.Body.Close(); err != nil {
				logger.Error(r.Context(), "failed decrypt with error", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			data, err := decryption.Decrypt(bodyBytes)
			if err != nil {
				logger.Error(r.Context(), "failed decrypt with error", err)
				w.WriteHeader(http.StatusBadRequest)
			}
			r.Body = io.NopCloser(bytes.NewBuffer(data))
		}

		next.ServeHTTP(w, r)
	})
}
