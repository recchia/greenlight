package data

import (
	"testing"
)

func TestRuntimeMarshalJSON(t *testing.T) {
	r := Runtime(102)
	b, err := r.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != `"102 mins"` {
		t.Errorf("expected %q, got %q", `"102 mins"`, string(b))
	}
}

func TestRuntimeUnmarshalJSON(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		var r Runtime
		err := r.UnmarshalJSON([]byte(`"102 mins"`))
		if err != nil {
			t.Fatal(err)
		}

		if r != 102 {
			t.Errorf("expected 102, got %d", r)
		}
	})

	t.Run("Invalid format", func(t *testing.T) {
		var r Runtime
		err := r.UnmarshalJSON([]byte(`"102"`))
		if err == nil {
			t.Error("expected error for invalid format")
		}
	})
}
