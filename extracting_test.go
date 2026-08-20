package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlrgo/result"
)

func TestResult_Expect(t *testing.T) {
	const msg = "expected result to be ok"
	tests := []struct {
		name      string
		give      result.Result[int, error]
		want      int
		wantPanic bool
	}{
		{"err", result.Err[int](ErrTest), 0, true},
		{"ok", result.Ok[int, error](10), 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.PanicsWithValue(t, fmt.Sprintf("%s: %v", msg, ErrTest), func() { tt.give.Expect(msg) })
				return
			}
			got := tt.give.Expect(msg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResult_ExpectErr(t *testing.T) {
	const msg = "expected result to be err"
	tests := []struct {
		name      string
		give      result.Result[int, error]
		want      error
		wantPanic any
	}{
		{"err", result.Err[int](ErrTest), ErrTest, nil},
		{"ok", result.Ok[int, error](10), nil, fmt.Sprintf("%s: %v", msg, 10)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic != nil {
				assert.PanicsWithValue(t, tt.wantPanic, func() { tt.give.ExpectErr(msg) })
				return
			}
			got := tt.give.ExpectErr(msg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResult_Unwrap(t *testing.T) {
	tests := []struct {
		name      string
		give      result.Result[int, error]
		want      int
		wantPanic any
	}{
		{"err", result.Err[int](ErrTest), 0, "result: Unwrap called on Err: test"},
		{"ok", result.Ok[int, error](10), 10, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic != nil {
				assert.PanicsWithValue(t, tt.wantPanic, func() { tt.give.Unwrap() })
				return
			}
			got := tt.give.Unwrap()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResult_UnwrapErr(t *testing.T) {
	tests := []struct {
		name      string
		give      result.Result[int, error]
		want      error
		wantPanic any
	}{
		{"err", result.Err[int](ErrTest), ErrTest, nil},
		{"ok", result.Ok[int, error](10), nil, "result: UnwrapErr called on Ok: 10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic != nil {
				assert.PanicsWithValue(t, tt.wantPanic, func() { tt.give.UnwrapErr() })
				return
			}
			got := tt.give.UnwrapErr()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResult_UnwrapOr(t *testing.T) {
	tests := []struct {
		name string
		give result.Result[int, error]
		want int
	}{
		{"err", result.Err[int](ErrTest), 67},
		{"ok", result.Ok[int, error](10), 10},
		{"ok zero", result.Ok[int, error](0), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.UnwrapOr(67)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResult_UnwrapOrDefault(t *testing.T) {
	tests := []struct {
		name string
		give result.Result[int, error]
		want int
	}{
		{"err", result.Err[int](ErrTest), 0},
		{"ok", result.Ok[int, error](10), 10},
		{"ok zero", result.Ok[int, error](0), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.UnwrapOrDefault()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResult_UnwrapOrElse(t *testing.T) {
	tests := []struct {
		name      string
		give      result.Result[int, error]
		want      int
		wantCalls int
	}{
		{"err", result.Err[int](ErrTest), 67, 1},
		{"ok", result.Ok[int, error](10), 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.give.UnwrapOrElse(func(err error) int {
				calls++
				assert.Equal(t, ErrTest, err)
				return 67
			})
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantCalls, calls)
		})
	}
}

func ExampleResult_Expect() {
	defer func() { fmt.Println(recover()) }()

	fmt.Println(result.Ok[string, error]("value").Expect("fruits are healthy"))
	result.Err[string](errors.New("empty")).Expect("string should not be empty")

	// Output:
	// value
	// string should not be empty: empty
}

func ExampleResult_ExpectErr() {
	defer func() { fmt.Println(recover()) }()

	fmt.Println(result.Err[int](errors.New("late")).ExpectErr("should fail"))
	result.Ok[int, error](10).ExpectErr("should fail")

	// Output:
	// late
	// should fail: 10
}

func ExampleResult_Unwrap() {
	defer func() { fmt.Println(recover()) }()

	fmt.Println(result.Ok[string, error]("air").Unwrap())
	result.Err[string](errors.New("late")).Unwrap()

	// Output:
	// air
	// result: Unwrap called on Err: late
}

func ExampleResult_UnwrapErr() {
	defer func() { fmt.Println(recover()) }()

	fmt.Println(result.Err[int](errors.New("late")).UnwrapErr())
	result.Ok[int, error](10).UnwrapErr()

	// Output:
	// late
	// result: UnwrapErr called on Ok: 10
}

func ExampleResult_UnwrapOr() {
	fmt.Println(result.Err[string](errors.New("late")).UnwrapOr("bike"))
	fmt.Println(result.Ok[string, error]("car").UnwrapOr("bike"))

	// Output:
	// bike
	// car
}

func ExampleResult_UnwrapOrDefault() {
	fmt.Println(result.Err[int](errors.New("late")).UnwrapOrDefault())
	fmt.Println(result.Ok[int, error](123).UnwrapOrDefault())

	// Output:
	// 0
	// 123
}

func ExampleResult_UnwrapOrElse() {
	fmt.Println(result.Err[int](errors.New("late")).UnwrapOrElse(func(error) int { return 10 }))
	fmt.Println(result.Ok[int, error](6).UnwrapOrElse(func(error) int { return 10 }))

	// Output:
	// 10
	// 6
}
