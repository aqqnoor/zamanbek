package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

type ctxKey string

const (
	ctxKeyReqID ctxKey = "req_id"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func genReqID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "req-unknown"
	}
	return hex.EncodeToString(buf)
}

// withRequestID: кладём X-Request-ID в контекст и в заголовок ответа
func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = genReqID()
		}
		ctx := context.WithValue(r.Context(), ctxKeyReqID, reqID)
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// withLogging: логируем метод, путь, код, длительность и req_id
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lrw, r)
		dur := time.Since(start)
		reqID := GetReqID(r)

		log.Printf("[http] %s %s status=%d bytes=%d dur_ms=%.1f req_id=%s",
			r.Method, r.URL.Path, lrw.status, lrw.bytes, float64(dur.Milliseconds()), reqID)
	})
}

// GetReqID: достаём значение из стандартного контекста
func GetReqID(r *http.Request) string {
	if v := r.Context().Value(ctxKeyReqID); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}