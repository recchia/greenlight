package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestShowMovieHandler(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Valid ID", func(t *testing.T) {
		code, _, body := ts.get(t, "/v1/movies/1")

		if code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, code)
		}

		var response struct {
			Movie struct {
				ID    int64  `json:"id"`
				Title string `json:"title"`
			} `json:"movie"`
		}

		err := json.Unmarshal([]byte(body), &response)
		if err != nil {
			t.Fatal(err)
		}

		if response.Movie.ID != 1 {
			t.Errorf("expected movie ID 1, got %d", response.Movie.ID)
		}
		if response.Movie.Title != "Test Movie" {
			t.Errorf("expected movie title 'Test Movie', got %q", response.Movie.Title)
		}
	})

	t.Run("Invalid ID", func(t *testing.T) {
		code, _, _ := ts.get(t, "/v1/movies/abc")

		if code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, code)
		}
	})

	t.Run("Non-existent ID", func(t *testing.T) {
		code, _, _ := ts.get(t, "/v1/movies/0")

		if code != http.StatusNotFound {
			t.Errorf("expected status code %d, got %d", http.StatusNotFound, code)
		}
	})
}
