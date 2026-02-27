package main

import (
	"net/http"
	"testing"
)

func TestRegisterUserHandler(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Valid registration", func(t *testing.T) {
		input := map[string]any{
			"name":     "Alice",
			"email":    "alice@example.com",
			"password": "pa$$word123",
		}

		code, _, _ := ts.postJSON(t, "/v1/users", input)

		if code != http.StatusAccepted {
			t.Errorf("expected status code %d, got %d", http.StatusAccepted, code)
		}
	})

	t.Run("Invalid email", func(t *testing.T) {
		input := map[string]any{
			"name":     "Alice",
			"email":    "invalid-email",
			"password": "pa$$word123",
		}

		code, _, _ := ts.postJSON(t, "/v1/users", input)

		if code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, got %d", http.StatusUnprocessableEntity, code)
		}
	})

	t.Run("Missing password", func(t *testing.T) {
		input := map[string]any{
			"name":  "Alice",
			"email": "alice@example.com",
		}

		code, _, _ := ts.postJSON(t, "/v1/users", input)

		if code != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, got %d", http.StatusUnprocessableEntity, code)
		}
	})
}
