package main

import (
	"net/http"
	"strconv"
	"time"
)

type statusTrackingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (w *statusTrackingResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusTrackingResponseWriter) Write(data []byte) (int, error) {
	size, err := w.ResponseWriter.Write(data)
	w.size += size
	return size, err
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/health" {
			next.ServeHTTP(w, r)
			return
		}
		startTime := time.Now()

		sw := &statusTrackingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(sw, r)

		elapsed := time.Since(startTime)

		if sw.statusCode == http.StatusOK {
			logger.Info("request", "method", r.Method, "path", r.URL.Path, "from", r.RemoteAddr, "status", strconv.Itoa(sw.statusCode), "size", strconv.Itoa(sw.size), "in", elapsed)
			return
		} else if sw.statusCode == http.StatusNotFound {
			logger.Warn("request", "method", r.Method, "path", r.URL.Path, "from", r.RemoteAddr, "status", strconv.Itoa(sw.statusCode), "size", strconv.Itoa(sw.size), "in", elapsed)
			return
		} else if sw.statusCode == http.StatusInternalServerError {
			logger.Error("request", "method", r.Method, "path", r.URL.Path, "from", r.RemoteAddr, "status", strconv.Itoa(sw.statusCode), "size", strconv.Itoa(sw.size), "in", elapsed)
			return
		} else {
			logger.Info("request", "method", r.Method, "path", r.URL.Path, "from", r.RemoteAddr, "status", strconv.Itoa(sw.statusCode), "size", strconv.Itoa(sw.size), "in", elapsed)
			return
		}

	})

}
