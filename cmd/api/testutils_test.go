package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/recchia/greenlight/internal/data"
	"github.com/recchia/greenlight/internal/mailer"
)

func newTestApplication(t *testing.T) *application {
	return &application{
		config: config{
			limiter: struct {
				rps     float64
				burst   int
				enabled bool
			}{
				enabled: false,
			},
		},
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		models: data.NewMockModels(),
		mailer: &mailer.Mailer{},
		wg:     sync.WaitGroup{},
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	req, err := http.NewRequest(http.MethodGet, ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add a dummy token to bypass authenticate middleware
	req.Header.Set("Authorization", "Bearer token26charslong1234567890")

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) postJSON(t *testing.T, urlPath string, body any) (int, http.Header, string) {
	js, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, ts.URL+urlPath, bytes.NewReader(js))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	resBody, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(resBody)
}
