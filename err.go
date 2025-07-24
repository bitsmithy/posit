package posit

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func Err(tb testing.TB, got error) {
	tb.Helper()

	if got == nil {
		tb.Error("want any error, got <nil>")
	}
}

func NoErr(tb testing.TB, got error) {
	tb.Helper()

	if got != nil {
		tb.Fatalf("want no error, got %T(%v)", got, got)
	}
}

func ErrIs(tb testing.TB, got error, want any) {
	tb.Helper()

	switch w := want.(type) {
	case error:
		if !errors.Is(got, w) {
			tb.Errorf("want %T(%v), got %T(%v)", w, w, got, got)
		}
	case reflect.Type:
		if !errors.As(got, &w) {
			tb.Errorf("want error type %T, got %T", w, got)
		}
	case string:
		if !strings.Contains(got.Error(), w) {
			tb.Errorf("want %q, got %q", w, got.Error())
		}
	case nil:
		if got != nil {
			tb.Fatalf("want <nil>, got %T(%v)", got, got)
			return
		}
	default:
		tb.Fatalf("did not want an error or an error type.")
	}
}
