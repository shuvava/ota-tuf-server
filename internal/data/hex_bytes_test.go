package data_test

import (
	"encoding/json"
	"testing"

	"github.com/shuvava/ota-tuf-server/internal/data"
)

func TestHexBytes(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		got, err := json.Marshal(data.HexBytes("foo"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		want := []byte(`"666f6f"`)
		if string(got) != string(want) {
			t.Errorf("expected %s, got %s", want, got)
		}
	})
	t.Run("UnmarshalJSON", func(t *testing.T) {
		var got data.HexBytes
		err := json.Unmarshal([]byte(`"666f6f"`), &got)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		want := "foo"
		if string(got) != want {
			t.Errorf("expected %s, got %s", want, got)
		}
	})
	t.Run("UnmarshalJSON error uneven length", func(t *testing.T) {
		var got data.HexBytes
		err := json.Unmarshal([]byte(`"a"`), &got)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
	t.Run("UnmarshalJSON error invalid hex", func(t *testing.T) {
		var got data.HexBytes
		err := json.Unmarshal([]byte(`"zz"`), &got)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
	t.Run("UnmarshalJSON error wrong type", func(t *testing.T) {
		var got data.HexBytes
		err := json.Unmarshal([]byte(`"6"`), &got)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
