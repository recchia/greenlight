package data

import (
	"testing"

	"github.com/recchia/greenlight/internal/validator"
)

func TestValidateMovie(t *testing.T) {
	v := validator.New()

	movie := &Movie{
		Title:   "Test Movie",
		Year:    2024,
		Runtime: 120,
		Genres:  []string{"action"},
	}

	ValidateMovie(v, movie)
	if !v.Valid() {
		t.Errorf("expected valid movie, got errors: %v", v.Errors)
	}

	t.Run("Missing title", func(t *testing.T) {
		v := validator.New()
		movie.Title = ""
		ValidateMovie(v, movie)
		if v.Valid() {
			t.Error("expected invalid movie due to missing title")
		}
		movie.Title = "Test Movie" // reset
	})

	t.Run("Future year", func(t *testing.T) {
		v := validator.New()
		movie.Year = 3000
		ValidateMovie(v, movie)
		if v.Valid() {
			t.Error("expected invalid movie due to future year")
		}
		movie.Year = 2024 // reset
	})
}
