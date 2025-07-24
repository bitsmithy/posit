package posit_test

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/bitsmithy/posit"
)

type wrappedString struct {
	val string
}

type implementsEqual struct {
	// This implements an Equal method that only checks that `val` matches.
	val   int
	other int
}

func newImplementsEqual(val int) implementsEqual {
	return implementsEqual{val: val, other: rand.Int()}
}

func (n implementsEqual) Equal(other implementsEqual) bool {
	return n.val == other.val
}

func TestEqual(t *testing.T) {
	t.Run("calls Helper", func(t *testing.T) {
		tb := &mockTB{}
		posit.Equal(tb, "foo", "bar")

		if !tb.helper {
			t.Error("should have called Helper() but did not")
		}
	})

	t.Run("equal", func(t *testing.T) {
		now := time.Now()
		ptr := "foo"

		testCases := map[string]struct {
			got  any
			want any
		}{
			"integer":     {got: 3, want: 3},
			"string":      {got: "boo", want: "boo"},
			"bool":        {got: true, want: true},
			"byte":        {got: byte(5), want: byte(5)},
			"rune":        {got: 'C', want: 'C'},
			"struct":      {got: wrappedString{"foo"}, want: wrappedString{"foo"}},
			"pointer":     {got: &ptr, want: &ptr},
			"nil slice":   {got: []int(nil), want: []int(nil)},
			"byte slice":  {got: []byte("foo"), want: []byte("foo")},
			"int slice":   {got: []int{42, 84}, want: []int{42, 84}},
			"time.Time":   {got: now, want: now},
			"nil":         {got: nil, want: nil},
			"nil pointer": {got: (*int)(nil), want: (*int)(nil)},
			"nil map":     {got: map[string]int(nil), want: map[string]int(nil)},
			"nil chan":    {got: (chan int)(nil), want: (chan int)(nil)},
			"empty map":   {got: map[string]int{}, want: map[string]int{}},
			"map":         {got: map[string]int{"foo": 1}, want: map[string]int{"foo": 1}},
		}

		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				tb := &mockTB{}
				posit.Equal(tb, tc.got, tc.want)
				if tb.failed {
					t.Errorf("%#v vs %#v: should be equal", tc.got, tc.want)
				}
			})
		}
	})

	t.Run("not equal", func(t *testing.T) {
		now := time.Now()
		foo := "foo"
		bar := "bar"

		testCases := map[string]struct {
			got  any
			want any
			msg  string
		}{
			"integer": {
				got: 1, want: 2,
				msg: "want 2, got 1",
			},
			"int32 vs int64": {
				got: int32(2), want: int64(2),
				msg: "want 2, got 2",
			},
			"int vs string": {
				got: 7, want: "7",
				msg: `want "7", got 7`,
			},
			"string": {
				got: "foo", want: "bar",
				msg: `want "bar", got "foo"`,
			},
			"bool": {
				got: true, want: false,
				msg: "want false, got true",
			},
			"struct": {
				got: wrappedString{"foo"}, want: wrappedString{"bar"},
				msg: `want posit_test.wrappedString{val:"bar"}, got posit_test.wrappedString{val:"foo"}`,
			},
			"pointer": {
				got: &foo, want: &bar,
			},
			"byte slice": {
				got: []byte("abc"), want: []byte("abd"),
				msg: `want []byte{0x61, 0x62, 0x64}, got []byte{0x61, 0x62, 0x63}`,
			},
			"int slice": {
				got: []int{33, 11}, want: []int{11, 33},
				msg: `want []int{11, 33}, got []int{33, 11}`,
			},
			"int slice vs any slice": {
				got: []int{12, 34}, want: []any{12, 34},
				msg: `want []interface {}{12, 34}, got []int{12, 34}`,
			},
			"time.Time": {
				got: now, want: now.Add(time.Second),
			},
			"nil vs non-nil": {
				got: nil, want: 1.2,
				msg: "want 1.2, got <nil>",
			},
			"non-nil vs nil": {
				got: 4, want: nil,
				msg: "want <nil>, got 4",
			},
			"nil vs empty": {
				got: []int(nil), want: []int{},
				msg: "want []int{}, got []int(nil)",
			},
			"map": {
				got: map[string]int{"foo": 2}, want: map[string]int{"foo": 1},
				msg: `want map[string]int{"foo":1}, got map[string]int{"foo":2}`,
			},
			"chan": {
				got: make(chan int), want: make(chan int),
			},
		}

		for name, tc := range testCases {
			t.Run(name, func(t *testing.T) {
				tb := &mockTB{}
				posit.Equal(tb, tc.got, tc.want)
				if !tb.failed {
					t.Errorf("%#v vs %#v: should have failed", tc.got, tc.want)
				}
				if tb.fatal {
					t.Error("should not be fatal")
				}
				if tc.msg != "" && tb.msg != tc.msg {
					t.Errorf("expected error message '%s', got '%s'", tc.msg, tb.msg)
				}
			})
		}
	})
	t.Run("time", func(t *testing.T) {
		// date1 and date2 represent the same point in time,
		date1 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		date2 := time.Date(2025, 1, 1, 5, 0, 0, 0, time.FixedZone("UTC+5", 5*3600))
		tb := &mockTB{}
		posit.Equal(tb, date1, date2)
		if tb.failed {
			t.Errorf("%#v vs %#v: should have passed", date1, date2)
		}
	})
	t.Run("equaler", func(t *testing.T) {
		t.Run("equal", func(t *testing.T) {
			tb := &mockTB{}
			n1, n2 := newImplementsEqual(4), newImplementsEqual(4)
			posit.Equal(tb, n1, n2)
			if tb.failed {
				t.Errorf("%#v vs %#v: should be equal", n1, n2)
			}
		})
		t.Run("non-equal", func(t *testing.T) {
			tb := &mockTB{}
			n1, n2 := newImplementsEqual(1), newImplementsEqual(2)
			posit.Equal(tb, n1, n2)
			if !tb.failed {
				t.Errorf("%#v vs %#v: should not be equal", n1, n2)
			}
			if tb.fatal {
				t.Error("should not be fatal")
			}
		})
	})
}

func TestEqualAny(t *testing.T) {
	t.Run("calls Helper", func(t *testing.T) {
		tb := &mockTB{}
		posit.EqualAny(tb, "foo")

		if !tb.helper {
			t.Error("should have called Helper() but did not")
		}
	})

	t.Run("no wants", func(t *testing.T) {
		tb := &mockTB{}
		posit.EqualAny(tb, "foo")

		if !tb.fatal {
			t.Error("should be fatal, as no wants were passed")
		}
	})

	t.Run("matches", func(t *testing.T) {
		t.Run("multiple wants", func(t *testing.T) {
			tb := &mockTB{}
			got := 7
			wants := []int{1, 3, 5, 7, 9}

			posit.EqualAny(tb, got, wants...)
			if tb.failed {
				t.Errorf("%#v should be in %#v", got, wants)
			}
		})

		t.Run("single wants", func(t *testing.T) {
			tb := &mockTB{}
			got := "foo"
			wants := []string{"foo"}
			posit.EqualAny(tb, got, wants...)

			if tb.failed {
				t.Errorf("%#v should be equal to %#v", got, wants)
			}
		})
	})

	t.Run("no matches", func(t *testing.T) {
		t.Run("multiple wants", func(t *testing.T) {
			tb := &mockTB{}
			got := 4
			wants := []int{1, 3, 5, 7, 9}
			errMsg := "want any of [1 3 5 7 9], got 4"

			posit.EqualAny(tb, got, wants...)
			if !tb.failed {
				t.Errorf("%#v should not be in %#v", got, wants)
			}
			if tb.fatal {
				t.Error("should not be fatal")
			}
			if errMsg != "" && tb.msg != errMsg {
				t.Errorf("expected error message '%s', got '%s'", errMsg, tb.msg)
			}
		})

		t.Run("single wants", func(t *testing.T) {
			tb := &mockTB{}
			got := "foo"
			wants := []string{"bar"}
			errMsg := `want "bar", got "foo"`

			posit.EqualAny(tb, got, wants...)
			if !tb.failed {
				t.Errorf("%#v should not be in %#v", got, wants)
			}
			if tb.fatal {
				t.Error("should not be fatal")
			}
			if errMsg != "" && tb.msg != errMsg {
				t.Errorf("expected error message '%s', got '%s'", errMsg, tb.msg)
			}
		})
	})
}
