package posit

import (
	"bytes"
	"reflect"
	"testing"
)

// equaler is an interface for types with an Equal method
// (like time.Time or net.IP).
type equaler[T any] interface {
	Equal(T) bool
}

func Equal[T any](tb testing.TB, got T, want T) {
	tb.Helper()

	if equal(got, want) {
		return
	}

	tb.Errorf("want %#v, got %#v", want, got)
}

func EqualAny[T any](tb testing.TB, got T, wants ...T) {
	tb.Helper()

	switch len(wants) {
	case 0:
		tb.Fatal("no wants given")
	case 1:
		Equal(tb, got, wants[0])
	default:
		for _, want := range wants {
			if equal(got, want) {
				return
			}
		}

		tb.Errorf("want any of %v, got %#v", wants, got)
	}
}

func equal[T any](a, b T) bool {
	// Check if both are nil.
	if isNil(a) && isNil(b) {
		return true
	}

	// use Equal method, if it exists
	if eq, ok := any(a).(equaler[T]); ok {
		return eq.Equal(b)
	}

	// use Equal for byte slices as well
	if aBytes, ok := any(a).([]byte); ok {
		bBytes := any(b).([]byte)
		return bytes.Equal(aBytes, bBytes)
	}

	// Fallback to reflective comparison.
	return reflect.DeepEqual(a, b)
}

func isNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return rv.IsNil()
	default:
		return false
	}
}
