package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/v1/healthcheck")

	if code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, code)
	}

	var response struct {
		Status     string `json:"status"`
		SystemInfo struct {
			Environment string `json:"environment"`
			Version     string `json:"version"`
		} `json:"system_info"`
	}

	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != "available" {
		t.Errorf("expected status 'available', got %q", response.Status)
	}
}
