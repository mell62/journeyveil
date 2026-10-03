package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

const requestIDHeader = "X-Request-ID"

type requestIDContextKey struct{}

type responseWriter struct {
	http.ResponseWriter
	status      int
	bytes       int
	wroteHeader bool
}

func (w *responseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Write(body []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	written, err := w.ResponseWriter.Write(body)
	w.bytes += written
	return written, err
}

func (w *responseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// RequestID returns the server-generated request ID stored in ctx.
func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDContextKey{}).(string)
	return requestID
}

func traceRequests(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, err := newRequestID()
		if err != nil {
			logger.ErrorContext(r.Context(), "generate request ID", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		started := time.Now()
		w.Header().Set(requestIDHeader, requestID)
		tracked := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		ctx := context.WithValue(r.Context(), requestIDContextKey{}, requestID)

		next.ServeHTTP(tracked, r.WithContext(ctx))

		logger.InfoContext(
			ctx,
			"HTTP request completed",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"status", tracked.status,
			"response_bytes", tracked.bytes,
			"duration_ms", time.Since(started).Milliseconds(),
		)
	})
}

func newRequestID() (string, error) {
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return "", err
	}

	id[6] = id[6]&0x0f | 0x40
	id[8] = id[8]&0x3f | 0x80

	var encoded [36]byte
	hex.Encode(encoded[0:8], id[0:4])
	encoded[8] = '-'
	hex.Encode(encoded[9:13], id[4:6])
	encoded[13] = '-'
	hex.Encode(encoded[14:18], id[6:8])
	encoded[18] = '-'
	hex.Encode(encoded[19:23], id[8:10])
	encoded[23] = '-'
	hex.Encode(encoded[24:36], id[10:16])

	return string(encoded[:]), nil
}
