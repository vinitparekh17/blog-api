package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusTrackingResponseWriter_Write(t *testing.T) {
	rr := httptest.NewRecorder()
	w := &statusTrackingResponseWriter{
		ResponseWriter: rr,
	}

	data := []byte("Hello, World!")

	size, err := w.Write(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if size != len(data) {
		t.Errorf("want size %d; got %d", len(data), size)
	}

	if w.size != len(data) {
		t.Errorf("want size %d; got %d", len(data), w.size)
	}

	if rr.Body.String() != string(data) {
		t.Errorf("want body %q; got %q", string(data), rr.Body.String())
	}
}

func TestStatusTrackingResponseWriter_WriteHeader(t *testing.T) {
	rr := httptest.NewRecorder()
	w := &statusTrackingResponseWriter{
		ResponseWriter: rr,
	}

	code := http.StatusOK
	w.WriteHeader(code)

	if w.statusCode != code {
		t.Errorf("want status code %d; got %d", code, w.statusCode)
	}

	if rr.Code != code {
		t.Errorf("want response code %d; got %d", code, rr.Code)
	}
}
