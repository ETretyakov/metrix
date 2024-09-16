// Module "middlewares" holds functions that are used as net/http middlewares.
package middlewares

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"metrix/pkg/logger"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header - the function that retrieves header from the reponse writer.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write - the function that writes to the response.
func (c *compressWriter) Write(p []byte) (int, error) {
	if n, err := c.zw.Write(p); err != nil {
		return 0, errors.Wrapf(err, "failed to close writer")
	} else {
		return n, nil
	}
}

// Write - the function that writes header to the response.
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode <= http.StatusMultipleChoices {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close - the function that closes the response writer.
func (c *compressWriter) Close() error {
	if err := c.zw.Close(); err != nil {
		return errors.Wrapf(err, "failed to close writer")
	}
	return nil
}

// GzipMiddleware - the middleware function that enables compression for http communicztion.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer func() {
				if err := cw.Close(); err != nil {
					logger.Error(context.TODO(), "failed to close compress writer", err)
					return
				}
			}()
			w.Header().Set("Content-Encoding", "gzip")
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Error(context.TODO(), "failed to get compress reader", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer func() {
				if err := reader.Close(); err != nil {
					logger.Error(context.TODO(), "failed to close reader", err)
				}
			}()

			decompressed, err := io.ReadAll(reader)
			if err != nil {
				logger.Error(context.TODO(), "failed to decompress", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(decompressed))
		}

		next.ServeHTTP(ow, r)
	})
}
