package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlrgo/result"
)

func TestResult_IsErr(t *testing.T) {
	tests := []struct {
		name string
		give result.Result[int, error]
		want bool
	}{
		{"err", result.Err[int](ErrTest), true},
		{"ok", result.Ok[int, error](123), false},
		{"zero", result.Result[int, error]{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.IsErr()
			assert.Equal(t, tt.want, got, "want: %v, got: %v", tt.want, got)
		})
	}
}

func TestResult_IsErrAnd(t *testing.T) {
	tests := []struct {
		name string
		give result.Result[int, error]
		pred func(error) bool
		want bool
	}{
		{"ok", result.Ok[int, error](123), nil, false},
		{"zero with nil err", result.Result[int, error]{}, func(err error) bool { return err == nil }, true},
		{
			"err with false predicate",
			result.Err[int](ErrTest),
			func(err error) bool { return err.Error() != ErrTest.Error() },
			false,
		},
		{
			"err with true predicate",
			result.Err[int](ErrTest),
			func(err error) bool { return err.Error() == ErrTest.Error() },
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.IsErrAnd(tt.pred)
			assert.Equal(t, tt.want, got, "want: %v, got: %v", tt.want, got)
		})
	}
}

func TestResult_IsOk(t *testing.T) {
	tests := []struct {
		name string
		give result.Result[int, error]
		want bool
	}{
		{"err", result.Err[int](ErrTest), false},
		{"ok", result.Ok[int, error](123), true},
		{"zero", result.Result[int, error]{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.IsOk()
			assert.Equal(t, tt.want, got, "want: %v, got: %v", tt.want, got)
		})
	}
}

func TestResult_IsOkAnd(t *testing.T) {
	tests := []struct {
		name string
		give result.Result[int, error]
		pred func(int) bool
		want bool
	}{
		{"err", result.Err[int](ErrTest), nil, false},
		{"zero", result.Result[int, error]{}, nil, false},
		{
			"ok with false predicate",
			result.Ok[int, error](123),
			func(i int) bool { return i != 123 },
			false,
		},
		{
			"ok with true predicate",
			result.Ok[int, error](123),
			func(i int) bool { return i == 123 },
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.IsOkAnd(tt.pred)
			assert.Equal(t, tt.want, got, "want: %v, got: %v", tt.want, got)
		})
	}
}

func ExampleResult_IsErr() {
	fmt.Println(result.Err[int](errors.New("late")).IsErr())
	fmt.Println(result.Ok[int, error](2).IsErr())

	// Output:
	// true
	// false
}

func ExampleResult_IsOk() {
	fmt.Println(result.Err[int](errors.New("late")).IsOk())
	fmt.Println(result.Ok[int, error](2).IsOk())

	// Output:
	// false
	// true
}

func ExampleResult_IsErrAnd() {
	isLate := func(err error) bool { return err.Error() == "late" }

	fmt.Println(result.Ok[int, error](2).IsErrAnd(isLate))
	fmt.Println(result.Err[int](errors.New("late")).IsErrAnd(isLate))

	// Output:
	// false
	// true
}

func ExampleResult_IsOkAnd() {
	positive := func(v int) bool { return v > 0 }

	fmt.Println(result.Err[int](errors.New("late")).IsOkAnd(positive))
	fmt.Println(result.Ok[int, error](5).IsOkAnd(positive))

	// Output:
	// false
	// true
}
