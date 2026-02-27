package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverPanic(t *testing.T) {
	app := newTestApplication(t)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	app.recoverPanic(h).ServeHTTP(rr, r)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	if rr.Header().Get("Connection") != "close" {
		t.Errorf("expected Connection: close header, got %q", rr.Header().Get("Connection"))
	}
}

func TestEnableCORS(t *testing.T) {
	app := newTestApplication(t)
	app.config.cors.trustedOrigins = []string{"http://localhost:3000"}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	t.Run("Preflight", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		r.Header.Set("Origin", "http://localhost:3000")
		r.Header.Set("Access-Control-Request-Method", "PATCH")

		app.enableCORS(h).ServeHTTP(rr, r)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
		if rr.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
			t.Errorf("expected Access-Control-Allow-Origin: http://localhost:3000, got %q", rr.Header().Get("Access-Control-Allow-Origin"))
		}
	})

	t.Run("Actual request", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Origin", "http://localhost:3000")

		app.enableCORS(h).ServeHTTP(rr, r)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
		if rr.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
			t.Errorf("expected Access-Control-Allow-Origin: http://localhost:3000, got %q", rr.Header().Get("Access-Control-Allow-Origin"))
		}
	})
}
