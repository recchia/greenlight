package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/recchia/greenlight/internal/validator"
)

func TestReadString(t *testing.T) {
	app := &application{}
	qs := url.Values{}
	qs.Add("name", "John")

	if app.readString(qs, "name", "Default") != "John" {
		t.Error("expected 'John'")
	}

	if app.readString(qs, "missing", "Default") != "Default" {
		t.Error("expected 'Default'")
	}
}

func TestReadCSV(t *testing.T) {
	app := &application{}
	qs := url.Values{}
	qs.Add("genres", "action,drama")

	csv := app.readCSV(qs, "genres", []string{})
	if len(csv) != 2 || csv[0] != "action" || csv[1] != "drama" {
		t.Error("expected ['action', 'drama']")
	}

	csv = app.readCSV(qs, "missing", []string{"default"})
	if len(csv) != 1 || csv[0] != "default" {
		t.Error("expected ['default']")
	}
}

func TestReadInt(t *testing.T) {
	app := &application{}
	v := validator.New()
	qs := url.Values{}
	qs.Add("page", "10")
	qs.Add("invalid", "abc")

	if app.readInt(qs, "page", 1, v) != 10 {
		t.Error("expected 10")
	}
	if !v.Valid() {
		t.Error("expected validator to be valid")
	}

	if app.readInt(qs, "invalid", 1, v) != 1 {
		t.Error("expected 1")
	}
	if v.Valid() {
		t.Error("expected validator to be invalid")
	}
}

func TestReadJSON(t *testing.T) {
	app := &application{}

	t.Run("Valid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"name": "John"}`)))
		var input struct {
			Name string `json:"name"`
		}

		err := app.readJSON(w, r, &input)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if input.Name != "John" {
			t.Errorf("expected 'John', got %q", input.Name)
		}
	})

	t.Run("Invalid JSON syntax", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"name": "John"`)))
		var input struct {
			Name string `json:"name"`
		}

		err := app.readJSON(w, r, &input)
		if err == nil {
			t.Error("expected error for invalid JSON syntax")
		}
	})

	t.Run("Unknown field", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`{"name": "John", "age": 30}`)))
		var input struct {
			Name string `json:"name"`
		}

		err := app.readJSON(w, r, &input)
		if err == nil {
			t.Error("expected error for unknown field")
		}
	})
}
