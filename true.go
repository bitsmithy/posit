package posit

import (
	"testing"
)

func True(tb testing.TB, got bool) {
	tb.Helper()

	if !got {
		tb.Error("want true, got false")
	}
}

func False(tb testing.TB, got bool) {
	tb.Helper()

	if got {
		tb.Error("want false, got true")
	}
}
