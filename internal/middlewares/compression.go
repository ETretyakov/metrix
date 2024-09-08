// Module "middlewares" holds functions that are used as net/http middlewares.
package middlewares

import (
	"compress/gzip"
	"io"
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

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create reader")
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read - the function that reads compressed data.
func (c compressReader) Read(p []byte) (n int, err error) {
	if n, err := c.zr.Read(p); err != nil {
		return 0, errors.Wrapf(err, "failed to read")
	} else {
		return n, nil
	}
}

// Close - the function that closes the reader of compressed data.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return errors.Wrapf(err, "fialed to close compress reader")
	}

	if err := c.zr.Close(); err != nil {
		return errors.Wrapf(err, "fialed to close compress reader")
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
					return
				}
			}()
			w.Header().Set("Content-Encoding", "gzip")
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer func() {
				if err := cr.Close(); err != nil {
					return
				}
			}()
		}

		next.ServeHTTP(ow, r)
	})
}
