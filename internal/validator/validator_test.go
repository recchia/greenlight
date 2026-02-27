package validator

import (
	"regexp"
	"testing"
)

func TestValidator(t *testing.T) {
	v := New()

	if !v.Valid() {
		t.Error("expected validator to be valid initially")
	}

	v.AddError("test", "test error")
	if v.Valid() {
		t.Error("expected validator to be invalid after adding error")
	}

	if v.Errors["test"] != "test error" {
		t.Errorf("expected error message to be 'test error', got %q", v.Errors["test"])
	}

	v.AddError("test", "another error")
	if v.Errors["test"] != "test error" {
		t.Error("expected AddError to not overwrite existing error")
	}
}

func TestValidatorCheck(t *testing.T) {
	v := New()

	v.Check(true, "test1", "should not be added")
	if len(v.Errors) != 0 {
		t.Error("expected no errors after successful check")
	}

	v.Check(false, "test2", "should be added")
	if v.Errors["test2"] != "should be added" {
		t.Error("expected error to be added after failed check")
	}
}

func TestPermittedValues(t *testing.T) {
	if !PermittedValues(1, 1, 2, 3) {
		t.Error("expected 1 to be permitted in [1, 2, 3]")
	}
	if PermittedValues(4, 1, 2, 3) {
		t.Error("expected 4 to not be permitted in [1, 2, 3]")
	}
}

func TestMatches(t *testing.T) {
	rx := regexp.MustCompile("^[a-z]+$")
	if !Matches("abc", rx) {
		t.Error("expected 'abc' to match ^[a-z]+$")
	}
	if Matches("123", rx) {
		t.Error("expected '123' not to match ^[a-z]+$")
	}
}

func TestUnique(t *testing.T) {
	if !Unique([]int{1, 2, 3}) {
		t.Error("expected [1, 2, 3] to be unique")
	}
	if Unique([]int{1, 2, 1}) {
		t.Error("expected [1, 2, 1] not to be unique")
	}
}
