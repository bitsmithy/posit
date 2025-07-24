package posit_test

import (
	"fmt"
	"testing"
)

// mockTB is a mock implementation of testing.TB
// to capture test failures.
type mockTB struct {
	testing.TB
	failed bool
	fatal  bool
	helper bool
	msg    string
}

func (m *mockTB) Helper() {
	m.helper = true
}

func (m *mockTB) Fatal(args ...any) {
	m.fatal = true
	m.Error(args...)
}

func (m *mockTB) Fatalf(format string, args ...any) {
	m.fatal = true
	m.Errorf(format, args...)
}

func (m *mockTB) Error(args ...any) {
	m.failed = true
	m.msg = fmt.Sprint(args...)
}

func (m *mockTB) Errorf(format string, args ...any) {
	m.failed = true
	m.msg = fmt.Sprintf(format, args...)
}
