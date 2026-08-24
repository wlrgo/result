package result_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlrgo/result/v2"
)

var (
	ErrTest  = errors.New("test")
	ErrOther = errors.New("other")
)

type typedNilError struct{}

func (*typedNilError) Error() string { return "typed nil" }

func TestCollect(t *testing.T) {
	tests := []struct {
		name      string
		give      []result.Result[int, error]
		wantValue []int
		wantErr   error
	}{
		{"empty nil", nil, []int{}, nil},
		{"empty slice", []result.Result[int, error]{}, []int{}, nil},
		{
			"all ok",
			[]result.Result[int, error]{result.Ok[int, error](1), result.Ok[int, error](2)},
			[]int{1, 2},
			nil,
		},
		{
			"first err",
			[]result.Result[int, error]{result.Err[int](ErrTest), result.Ok[int, error](2)},
			nil,
			ErrTest,
		},
		{
			"middle err",
			[]result.Result[int, error]{result.Ok[int, error](1), result.Err[int](ErrOther), result.Ok[int, error](3)},
			nil,
			ErrOther,
		},
		{
			"last err",
			[]result.Result[int, error]{result.Ok[int, error](1), result.Err[int](ErrTest)},
			nil,
			ErrTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := result.Collect(tt.give)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, got.UnwrapErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

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

func TestFrom(t *testing.T) {
	tests := []struct {
		name      string
		v         int
		err       error
		wantValue int
		wantErr   error
	}{
		{"nil error", 7, nil, 7, nil},
		{"zero value nil error", 0, nil, 0, nil},
		{"error", 7, ErrTest, 0, ErrTest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := result.From(tt.v, tt.err)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, got.UnwrapErr())
				return
			}

			assert.True(t, got.IsOk())
			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}

	t.Run("typed nil error", func(t *testing.T) {
		var err error = (*typedNilError)(nil)
		got := result.From(1, err)
		assert.Equal(t, err, got.UnwrapErr())
	})
}

func TestResult_And(t *testing.T) {
	tests := []struct {
		name      string
		a         result.Result[int, error]
		b         result.Result[string, error]
		wantValue string
		wantErr   error
	}{
		{
			"ok err",
			result.Ok[int, error](2),
			result.Err[string](ErrTest),
			"",
			ErrTest,
		},
		{
			"err ok",
			result.Err[int](ErrTest),
			result.Ok[string, error]("2"),
			"",
			ErrTest,
		},
		{
			"ok ok",
			result.Ok[int, error](2),
			result.Ok[string, error]("2"),
			"2",
			nil,
		},
		{
			"err err",
			result.Err[int](ErrTest),
			result.Err[string](ErrOther),
			"",
			ErrTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.And(tt.b)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, got.UnwrapErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.UnwrapOr("<none>"))
		})
	}
}

func TestResult_AndThen(t *testing.T) {
	tests := []struct {
		name      string
		a         result.Result[int, error]
		thenErr   bool
		wantValue string
		wantErr   error
		wantCalls int
		wantArg   int
	}{
		{
			"ok",
			result.Ok[int, error](2),
			false,
			"4",
			nil,
			1,
			2,
		},
		{
			"ok to err",
			result.Ok[int, error](2),
			true,
			"",
			ErrOther,
			1,
			2,
		},
		{
			"err",
			result.Err[int](ErrTest),
			false,
			"",
			ErrTest,
			0,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.a.AndThen(func(i int) result.Result[string, error] {
				calls++
				assert.Equal(t, tt.wantArg, i)
				if tt.thenErr {
					return result.Err[string](ErrOther)
				}
				return result.Ok[string, error](strconv.Itoa(i * i))
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, got.UnwrapErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.UnwrapOr("<none>"))
		})
	}
}

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

func TestResult_Get(t *testing.T) {
	tests := []struct {
		name     string
		give     result.Result[int, error]
		wantT    int
		wantSome bool
	}{
		{"err", result.Err[int](ErrTest), 0, false},
		{"ok", result.Ok[int, error](7), 7, true},
		{"ok zero", result.Ok[int, error](0), 0, true},
		{"zero value", result.Result[int, error]{}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.give.Get()
			assert.Equal(t, tt.wantT, got)
			assert.Equal(t, tt.wantSome, ok)
		})
	}
}

func TestResult_GetErr(t *testing.T) {
	tests := []struct {
		name     string
		give     result.Result[int, error]
		wantE    error
		wantSome bool
	}{
		{"err", result.Err[int](ErrTest), ErrTest, true},
		{"ok", result.Ok[int, error](7), nil, false},
		{"ok zero", result.Ok[int, error](0), nil, false},
		{"zero value", result.Result[int, error]{}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.give.GetErr()
			assert.Equal(t, tt.wantE, got)
			assert.Equal(t, tt.wantSome, ok)
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

func TestResult_Map(t *testing.T) {
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
			got := tt.give.Map(func(v string) int {
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

func TestResult_MapErr(t *testing.T) {
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
			got := tt.give.MapErr(func(err error) string {
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

func TestResult_MapOr(t *testing.T) {
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
			got := tt.give.MapOr(42, func(v string) int {
				calls++
				return len(v)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantCalls, calls)
		})
	}
}

func TestResult_MapOrDefault(t *testing.T) {
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
			got := tt.give.MapOrDefault(func(v string) int {
				calls++
				return len(v)
			})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantCalls, calls)
		})
	}
}

func TestResult_MapOrElse(t *testing.T) {
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
			got := tt.give.MapOrElse(
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

func TestResult_Or(t *testing.T) {
	tests := []struct {
		name      string
		a, b      result.Result[int, error]
		wantValue int
		wantErr   error
	}{
		{
			"ok err",
			result.Ok[int, error](2),
			result.Err[int](ErrTest),
			2,
			nil,
		},
		{
			"err ok",
			result.Err[int](ErrTest),
			result.Ok[int, error](67),
			67,
			nil,
		},
		{
			"ok ok",
			result.Ok[int, error](2),
			result.Ok[int, error](67),
			2,
			nil,
		},
		{
			"err err",
			result.Err[int](ErrTest),
			result.Err[int](ErrOther),
			-1,
			ErrOther,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Or(tt.b)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, got.UnwrapErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.UnwrapOr(-1))
		})
	}
}

func TestResult_OrElse(t *testing.T) {
	tests := []struct {
		name      string
		give      result.Result[int, error]
		thenErr   bool
		wantValue int
		wantErr   error
		wantCalls int
	}{
		{
			"ok",
			result.Ok[int, error](2),
			false,
			2,
			nil,
			0,
		},
		{
			"err to ok",
			result.Err[int](ErrTest),
			false,
			67,
			nil,
			1,
		},
		{
			"err to err",
			result.Err[int](ErrTest),
			true,
			-1,
			ErrOther,
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := tt.give.OrElse(func(err error) result.Result[int, error] {
				calls++
				assert.Equal(t, ErrTest, err)
				if tt.thenErr {
					return result.Err[int](ErrOther)
				}
				return result.Ok[int, error](67)
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, got.UnwrapErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.UnwrapOr(-1))
		})
	}
}

func TestResult_Seq(t *testing.T) {
	tests := []struct {
		name string
		give result.Result[int, error]
		want []int
	}{
		{"err", result.Err[int](ErrTest), nil},
		{"ok", result.Ok[int, error](2), []int{2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []int
			for v := range tt.give.Seq() {
				got = append(got, v)
			}
			assert.Equal(t, tt.want, got)
		})
	}

	t.Run("ok break", func(t *testing.T) {
		n := 0
		for range result.Ok[int, error](2).Seq() {
			n++
			break
		}
		assert.Equal(t, 1, n)
	})
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

func TestUnpack(t *testing.T) {
	tests := []struct {
		name   string
		give   result.Result[int, error]
		wantT  int
		wantE  error
		wantOk bool
	}{
		{"err", result.Err[int](ErrTest), 0, ErrTest, false},
		{"ok", result.Ok[int, error](7), 7, nil, true},
		{"ok zero", result.Ok[int, error](0), 0, nil, true},
		{"nil err", result.Err[int](error(nil)), 0, nil, false},
		{"zero value", result.Result[int, error]{}, 0, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotE := result.Unpack(tt.give)
			assert.Equal(t, tt.wantT, gotT)
			assert.Equal(t, tt.wantE, gotE)
			assert.Equal(t, tt.wantOk, tt.give.IsOk())
		})
	}
}

func ExampleCollect() {
	fmt.Println(
		result.Collect([]result.Result[int, error]{result.Ok[int, error](1), result.Ok[int, error](2)}).
			UnwrapOr([]int{-1}),
	)
	fmt.Println(
		result.Collect([]result.Result[int, error]{result.Ok[int, error](1), result.Err[int](errors.New("late"))}).
			UnwrapOr([]int{-1}),
	)

	// Output:
	// [1 2]
	// [-1]
}

func ExampleErr() {
	res := result.Err[int](errors.New("late"))

	fmt.Println(res.IsErr())
	// Output: true
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

func ExampleFrom() {
	parse := func(s string) (int, error) {
		if s == "7" {
			return 7, nil
		}
		return 0, errors.New("invalid")
	}

	fmt.Println(result.From(parse("7")).UnwrapOr(-1))
	fmt.Println(result.From(parse("x")).UnwrapOr(-1))

	// Output:
	// 7
	// -1
}

func ExampleOk() {
	res := result.Ok[int, error](10)

	fmt.Println(res.IsOk())
	fmt.Println(res.Unwrap())
	// Output:
	// true
	// 10
}

func ExampleResult() {
	var res result.Result[int, error]

	fmt.Println(res.IsErr())
	// Output: true
}

func ExampleResult_And() {
	x := result.Ok[int, error](2)
	y := result.Err[string](errors.New("late"))
	fmt.Println(x.And(y).UnwrapOr("-"))

	x = result.Err[int](errors.New("early"))
	y = result.Ok[string, error]("foo")
	fmt.Println(x.And(y).UnwrapOr("-"))

	x = result.Ok[int, error](2)
	y = result.Ok[string, error]("foo")
	fmt.Println(x.And(y).UnwrapOr("-"))

	x = result.Err[int](errors.New("early"))
	y = result.Err[string](errors.New("late"))
	fmt.Println(x.And(y).UnwrapOr("-"))

	// Output:
	// -
	// -
	// foo
	// -
}

func ExampleResult_AndThen() {
	sqThenToString := func(x int) result.Result[string, error] {
		if x > 10_000 {
			return result.Err[string](errors.New("overflow"))
		}
		return result.Ok[string, error](strconv.Itoa(x * x))
	}

	fmt.Println(result.Ok[int, error](2).AndThen(sqThenToString).UnwrapOr("-"))
	fmt.Println(result.Ok[int, error](1_000_000).AndThen(sqThenToString).UnwrapOr("-"))
	fmt.Println(result.Err[int](errors.New("empty")).AndThen(sqThenToString).UnwrapOr("-"))

	// Output:
	// 4
	// -
	// -
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

func ExampleResult_Get() {
	if v, ok := result.Ok[int, error](7).Get(); ok {
		fmt.Println(v)
	}

	if _, ok := result.Err[int](errors.New("late")).Get(); !ok {
		fmt.Println("err")
	}

	// Output:
	// 7
	// err
}

func ExampleResult_GetErr() {
	if err, ok := result.Err[int](errors.New("late")).GetErr(); ok {
		fmt.Println(err)
	}

	if _, ok := result.Ok[int, error](7).GetErr(); !ok {
		fmt.Println("ok")
	}

	// Output:
	// late
	// ok
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

func ExampleResult_IsErr() {
	fmt.Println(result.Err[int](errors.New("late")).IsErr())
	fmt.Println(result.Ok[int, error](2).IsErr())

	// Output:
	// true
	// false
}

func ExampleResult_IsErrAnd() {
	isLate := func(err error) bool { return err.Error() == "late" }

	fmt.Println(result.Ok[int, error](2).IsErrAnd(isLate))
	fmt.Println(result.Err[int](errors.New("late")).IsErrAnd(isLate))

	// Output:
	// false
	// true
}

func ExampleResult_IsOk() {
	fmt.Println(result.Err[int](errors.New("late")).IsOk())
	fmt.Println(result.Ok[int, error](2).IsOk())

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

func ExampleResult_Map() {
	x := result.Ok[string, error]("Hello, World!")
	fmt.Println(x.Map(func(v string) int { return len(v) }).UnwrapOr(-1))

	y := result.Err[string](errors.New("late"))
	fmt.Println(y.Map(func(v string) int { return len(v) }).UnwrapOr(-1))

	// Output:
	// 13
	// -1
}

func ExampleResult_MapErr() {
	x := result.Ok[string, error]("foo")
	fmt.Println(x.MapErr(func(err error) string { return err.Error() }).UnwrapOr("-"))

	y := result.Err[string](errors.New("late"))
	fmt.Println(y.MapErr(func(err error) string { return err.Error() }).UnwrapErr())

	// Output:
	// foo
	// late
}

func ExampleResult_MapOr() {
	x := result.Ok[string, error]("foo")
	fmt.Println(x.MapOr(42, func(v string) int { return len(v) }))

	x = result.Err[string](errors.New("late"))
	fmt.Println(x.MapOr(42, func(v string) int { return len(v) }))

	// Output:
	// 3
	// 42
}

func ExampleResult_MapOrDefault() {
	x := result.Ok[string, error]("hi")
	y := result.Err[string](errors.New("late"))

	fmt.Println(x.MapOrDefault(func(v string) int { return len(v) }))
	fmt.Println(y.MapOrDefault(func(v string) int { return len(v) }))

	// Output:
	// 2
	// 0
}

func ExampleResult_MapOrElse() {
	i := 21

	x := result.Ok[string, error]("foo")
	fmt.Println(
		x.MapOrElse(
			func(error) int { return 2 * i },
			func(v string) int { return len(v) },
		),
	)

	x = result.Err[string](errors.New("late"))
	fmt.Println(
		x.MapOrElse(
			func(error) int { return 2 * i },
			func(v string) int { return len(v) },
		),
	)

	// Output:
	// 3
	// 42
}

func ExampleResult_Or() {
	x := result.Ok[int, error](2)
	y := result.Err[int](errors.New("late"))
	fmt.Println(x.Or(y).UnwrapOr(-1))

	x = result.Err[int](errors.New("early"))
	y = result.Ok[int, error](100)
	fmt.Println(x.Or(y).UnwrapOr(-1))

	x = result.Ok[int, error](2)
	y = result.Ok[int, error](100)
	fmt.Println(x.Or(y).UnwrapOr(-1))

	x = result.Err[int](errors.New("early"))
	y = result.Err[int](errors.New("late"))
	fmt.Println(x.Or(y).UnwrapOr(-1))

	// Output:
	// 2
	// 100
	// 2
	// -1
}

func ExampleResult_OrElse() {
	nobody := func(error) result.Result[string, error] {
		return result.Err[string](errors.New("nobody"))
	}
	vikings := func(error) result.Result[string, error] {
		return result.Ok[string, error]("vikings")
	}

	fmt.Println(result.Ok[string, error]("barbarians").OrElse(vikings).UnwrapOr("-"))
	fmt.Println(result.Err[string](errors.New("empty")).OrElse(vikings).UnwrapOr("-"))
	fmt.Println(result.Err[string](errors.New("empty")).OrElse(nobody).UnwrapOr("-"))

	// Output:
	// barbarians
	// vikings
	// -
}

func ExampleResult_Seq() {
	for v := range result.Ok[int, error](2).Seq() {
		fmt.Println(v)
	}

	for range result.Err[int](errors.New("late")).Seq() {
		fmt.Println("err")
	}

	// Output:
	// 2
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

func ExampleUnpack() {
	fmt.Println(result.Unpack(result.Ok[int, error](7)))
	fmt.Println(result.Unpack(result.Err[int](errors.New("late"))))

	// Output:
	// 7 <nil>
	// 0 late
}
