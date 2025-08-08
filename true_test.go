package posit_test

import (
	"testing"

	"github.com/bitsmithy/posit"
)

func TestTrue(t *testing.T) {
	t.Run("calls Helper", func(t *testing.T) {
		tb := &mockTB{}
		posit.True(tb, 2 < 1)

		if !tb.helper {
			t.Error("should have called Helper() but did not")
		}
	})

	t.Run("with false", func(t *testing.T) {
		tb := &mockTB{}
		posit.True(tb, 5 < 2)

		if !tb.failed {
			t.Error("should have failed, as 5 > 2")
		}

		errMsg := "want true, got false"
		if tb.msg != errMsg {
			t.Errorf("expected error message '%s', got '%s'", errMsg, tb.msg)
		}
	})

	t.Run("with true", func(t *testing.T) {
		tb := &mockTB{}
		posit.True(tb, 5 > 2)

		if tb.failed {
			t.Error("should not have failed, as 5 > 2")
		}
	})
}

func TestFalse(t *testing.T) {
	t.Run("calls Helper", func(t *testing.T) {
		tb := &mockTB{}
		posit.False(tb, 2 > 1)

		if !tb.helper {
			t.Error("should have called Helper() but did not")
		}
	})

	t.Run("with true", func(t *testing.T) {
		tb := &mockTB{}
		posit.False(tb, 5 > 2)

		if !tb.failed {
			t.Error("should have failed, as 5 > 2")
		}

		errMsg := "want false, got true"
		if tb.msg != errMsg {
			t.Errorf("expected error message '%s', got '%s'", errMsg, tb.msg)
		}
	})

	t.Run("with false", func(t *testing.T) {
		tb := &mockTB{}
		posit.False(tb, 5 < 2)

		if tb.failed {
			t.Error("should not have failed, as 5 > 2")
		}
	})
}
