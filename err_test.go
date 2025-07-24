package posit_test

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/bitsmithy/posit"
)

func TestErr(t *testing.T) {
	t.Run("calls Helper", func(t *testing.T) {
		tb := &mockTB{}
		posit.Err(tb, errors.New("foo"))

		if !tb.helper {
			t.Error("should have called Helper() but did not")
		}
	})

	t.Run("without error", func(t *testing.T) {
		tb := &mockTB{}
		posit.Err(tb, nil)

		if !tb.failed {
			t.Error("should have failed, as an error was not passed in")
		}

		errMsg := "want any error, got <nil>"
		if tb.msg != errMsg {
			t.Errorf("expected error message '%s', got '%s'", errMsg, tb.msg)
		}
	})

	t.Run("with error", func(t *testing.T) {
		tb := &mockTB{}
		posit.Err(tb, errors.New("foo"))

		if tb.failed {
			t.Error("should not have failed, as an error was passed in")
		}
	})
}

func TestNoErr(t *testing.T) {
	t.Run("calls Helper", func(t *testing.T) {
		tb := &mockTB{}
		posit.NoErr(tb, nil)

		if !tb.helper {
			t.Error("should have called Helper() but did not")
		}
	})

	t.Run("without error", func(t *testing.T) {
		tb := &mockTB{}
		posit.NoErr(tb, nil)

		if tb.failed {
			t.Error("should not have failed, as an error was not passed in")
		}
	})

	t.Run("with error", func(t *testing.T) {
		tb := &mockTB{}
		posit.NoErr(tb, errors.New("foo"))

		if !tb.fatal {
			t.Error("should have failed, as an error was passed in")
		}

		errMsg := `want no error, got *errors.errorString(foo)`
		if tb.msg != errMsg {
			t.Errorf("expected error message '%s', got '%s'", errMsg, tb.msg)
		}
	})
}

func TestErrIs(t *testing.T) {
	_, openErr := os.Open("non-existing")

	t.Run("calls Helper", func(t *testing.T) {
		tb := &mockTB{}
		posit.ErrIs(tb, errors.New("bar"), errors.New("bar"))

		if !tb.helper {
			t.Error("should have called Helper() but did not")
		}
	})

	t.Run("want and got same error", func(t *testing.T) {
		tb := &mockTB{}

		_, err := os.Open("non-existing")
		posit.ErrIs(tb, err, fs.ErrNotExist)

		if tb.failed {
			t.Errorf("expected to pass, as %T(%v) should be the same as %T(%v)", err, err, fs.ErrNotExist, fs.ErrNotExist)
		}
	})

	t.Run("want type and got error of type", func(t *testing.T) {
		tb := &mockTB{}

		_, err := os.Open("non-existing")
		posit.ErrIs(tb, err, fs.ErrNotExist)

		if tb.failed {
			t.Errorf("expected to pass, as %T(%v) should be of type %T", err, err, fs.ErrNotExist)
		}
	})

	t.Run("want and got different error", func(t *testing.T) {
		tb := &mockTB{}

		err := errors.New("foo")
		posit.ErrIs(tb, errors.New("foo"), openErr)

		if !tb.failed {
			t.Errorf("expected to fail, as %T(%v) should not be the same as %T(%v)", err, err, openErr, openErr)
		}
	})

	t.Run("want type and got error of wrong type", func(t *testing.T) {
		tb := &mockTB{}

		posit.ErrIs(tb, openErr, fs.ErrClosed)

		if !tb.failed {
			t.Errorf("expected to fail, as %T(%v) should be of type %T", openErr, openErr, fs.ErrClosed)
		}
	})

	t.Run("want string and got matching sub-string", func(t *testing.T) {
		tb := &mockTB{}

		err := errors.New("foo is not acceptable")
		posit.ErrIs(tb, err, "foo")

		if tb.failed {
			t.Errorf("expected to pass, as %T(%v) should contain the string %q", err, err, "foo")
		}
	})

	t.Run("want string and got no matching sub-strings", func(t *testing.T) {
		tb := &mockTB{}

		err := errors.New("foo is not acceptable")
		posit.ErrIs(tb, err, "bar")

		if !tb.failed {
			t.Errorf("expected to fail, as %T(%v) should not contain the string %q", err, err, "bar")
		}
	})

	t.Run("want nil but got error", func(t *testing.T) {
		tb := &mockTB{}
		posit.ErrIs(tb, openErr, nil)

		if !tb.fatal {
			t.Errorf("expected to fail fatally, as no error was expected, but %T(%v) occurred.", openErr, openErr)
		}

		errMsg := "want <nil>, got *fs.PathError(open non-existing: no such file or directory)"
		if tb.msg != errMsg {
			t.Errorf("expected error message '%s', got '%s'", errMsg, tb.msg)
		}
	})

	t.Run("want error but got nil", func(t *testing.T) {
		tb := &mockTB{}
		posit.ErrIs(tb, nil, openErr)

		if !tb.failed {
			t.Errorf("expected to fail as %T(%v) was expected, but no error occurred", openErr, openErr)
		}
	})

	t.Run("want non-error", func(t *testing.T) {
		tb := &mockTB{}
		posit.ErrIs(tb, openErr, 1)

		if !tb.fatal {
			t.Error("expected to fail fatally, as want was not an error or error type")
		}
	})
}
