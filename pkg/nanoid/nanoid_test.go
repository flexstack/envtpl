package nanoid

import (
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("default length", func(t *testing.T) {
		id := New(21)

		if len(id) != 21 {
			t.Errorf("expected %d, got %d", 21, len(id))
		}
	})

	t.Run("custom length", func(t *testing.T) {
		for i := 1; i < 1024; i++ {
			id := New(i)

			if len(id) != i {
				t.Errorf("expected %d, got %d", i, len(id))
			}
		}
	})

	t.Run("custom alphabet", func(t *testing.T) {
		id := New(24, AlphabetAscii85)

		if len(id) != 24 {
			t.Errorf("expected %d, got %d", 24, len(id))
		}
	})
}
