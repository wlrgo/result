package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlrgo/result"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		name      string
		give      result.Result[result.Result[int, error], error]
		wantValue int
		wantErr   error
	}{
		{"ok ok", result.Ok[result.Result[int, error], error](result.Ok[int, error](6)), 6, nil},
		{"ok err", result.Ok[result.Result[int, error], error](result.Err[int](ErrTest)), 0, ErrTest},
		{"err", result.Err[result.Result[int, error]](ErrOther), 0, ErrOther},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := result.Flatten(tt.give)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, got.UnwrapErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		name      string
		give      result.Result[string, error]
		wantValue int
		wantErr   bool
		wantCalls int
	}{
		{"ok", result.Ok[string, error]("foo"), 3, false, 1},
		{"err", result.Err[string](ErrTest), 0, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := result.Map(tt.give, func(v string) int {
				calls++
				return len(v)
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantErr {
				assert.Equal(t, ErrTest, got.UnwrapErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestMapErr(t *testing.T) {
	tests := []struct {
		name      string
		give      result.Result[string, error]
		wantValue string
		wantErr   bool
		wantCalls int
	}{
		{"ok", result.Ok[string, error]("foo"), "foo", false, 0},
		{"err", result.Err[string](ErrTest), "", true, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := result.MapErr(tt.give, func(err error) string {
				calls++
				assert.Equal(t, ErrTest, err)
				return err.Error()
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantErr {
				assert.Equal(t, ErrTest.Error(), got.UnwrapErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestMapOr(t *testing.T) {
	tests := []struct {
		name      string
		give      result.Result[string, error]
		want      int
		wantCalls int
	}{
		{"ok", result.Ok[string, error]("foo"), 3, 1},
		{"err", result.Err[string](ErrTest), 42, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := result.MapOr(tt.give, 42, func(v string) int {
				calls++
				return len(v)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantCalls, calls)
		})
	}
}

func TestMapOrDefault(t *testing.T) {
	tests := []struct {
		name      string
		give      result.Result[string, error]
		want      int
		wantCalls int
	}{
		{"ok", result.Ok[string, error]("hi"), 2, 1},
		{"err", result.Err[string](ErrTest), 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := result.MapOrDefault(tt.give, func(v string) int {
				calls++
				return len(v)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantCalls, calls)
		})
	}
}

func TestMapOrElse(t *testing.T) {
	tests := []struct {
		name             string
		give             result.Result[string, error]
		want             int
		wantMapCalls     int
		wantDefaultCalls int
	}{
		{"ok", result.Ok[string, error]("foo"), 3, 1, 0},
		{"err", result.Err[string](ErrTest), 42, 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapCalls := 0
			defaultCalls := 0
			got := result.MapOrElse(
				tt.give,
				func(err error) int {
					defaultCalls++
					assert.Equal(t, ErrTest, err)
					return 42
				},
				func(v string) int {
					mapCalls++
					return len(v)
				},
			)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantMapCalls, mapCalls)
			assert.Equal(t, tt.wantDefaultCalls, defaultCalls)
		})
	}
}

func TestResult_Inspect(t *testing.T) {
	tests := []struct {
		name      string
		give      result.Result[int, error]
		wantValue int
		wantErr   bool
		wantCalls int
	}{
		{"ok", result.Ok[int, error](2), 2, false, 1},
		{"err", result.Err[int](ErrTest), 0, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.give.Inspect(func(v int) {
				calls++
				assert.Equal(t, tt.wantValue, v)
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantErr {
				assert.Equal(t, ErrTest, got.UnwrapErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestResult_InspectErr(t *testing.T) {
	tests := []struct {
		name      string
		give      result.Result[int, error]
		wantValue int
		wantErr   error
		wantCalls int
	}{
		{"ok", result.Ok[int, error](2), 2, nil, 0},
		{"err", result.Err[int](ErrTest), 0, ErrTest, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.give.InspectErr(func(err error) {
				calls++
				assert.Equal(t, tt.wantErr, err)
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, got.UnwrapErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func ExampleFlatten() {
	x := result.Ok[result.Result[int, error], error](result.Ok[int, error](6))
	fmt.Println(result.Flatten(x).UnwrapOr(-1))

	x = result.Ok[result.Result[int, error], error](result.Err[int](errors.New("inner")))
	fmt.Println(result.Flatten(x).UnwrapOr(-1))

	x = result.Err[result.Result[int, error]](errors.New("outer"))
	fmt.Println(result.Flatten(x).UnwrapOr(-1))

	// Output:
	// 6
	// -1
	// -1
}

func ExampleMap() {
	x := result.Ok[string, error]("Hello, World!")
	fmt.Println(result.Map(x, func(v string) int { return len(v) }).UnwrapOr(-1))

	y := result.Err[string](errors.New("late"))
	fmt.Println(result.Map(y, func(v string) int { return len(v) }).UnwrapOr(-1))

	// Output:
	// 13
	// -1
}

func ExampleMapErr() {
	x := result.Ok[string, error]("foo")
	fmt.Println(result.MapErr(x, func(err error) string { return err.Error() }).UnwrapOr("-"))

	y := result.Err[string](errors.New("late"))
	fmt.Println(result.MapErr(y, func(err error) string { return err.Error() }).UnwrapErr())

	// Output:
	// foo
	// late
}

func ExampleMapOr() {
	x := result.Ok[string, error]("foo")
	fmt.Println(result.MapOr(x, 42, func(v string) int { return len(v) }))

	x = result.Err[string](errors.New("late"))
	fmt.Println(result.MapOr(x, 42, func(v string) int { return len(v) }))

	// Output:
	// 3
	// 42
}

func ExampleMapOrDefault() {
	x := result.Ok[string, error]("hi")
	y := result.Err[string](errors.New("late"))

	fmt.Println(result.MapOrDefault(x, func(v string) int { return len(v) }))
	fmt.Println(result.MapOrDefault(y, func(v string) int { return len(v) }))

	// Output:
	// 2
	// 0
}

func ExampleMapOrElse() {
	i := 21

	x := result.Ok[string, error]("foo")
	fmt.Println(
		result.MapOrElse(
			x,
			func(error) int { return 2 * i },
			func(v string) int { return len(v) },
		),
	)

	x = result.Err[string](errors.New("late"))
	fmt.Println(
		result.MapOrElse(
			x,
			func(error) int { return 2 * i },
			func(v string) int { return len(v) },
		),
	)

	// Output:
	// 3
	// 42
}

func ExampleResult_Inspect() {
	x := result.Ok[int, error](2).Inspect(func(v int) { fmt.Println("got:", v) })

	fmt.Println(x.Unwrap())

	result.Err[int](errors.New("late")).Inspect(func(v int) { fmt.Println("got:", v) })

	// Output:
	// got: 2
	// 2
}

func ExampleResult_InspectErr() {
	result.Ok[int, error](2).InspectErr(func(err error) { fmt.Println("err:", err) })

	x := result.Err[int](errors.New("late")).InspectErr(func(err error) { fmt.Println("err:", err) })
	fmt.Println(x.IsErr())

	// Output:
	// err: late
	// true
}
