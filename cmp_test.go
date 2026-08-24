package result_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlrgo/result/v2"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		name string
		a, b result.Result[int, string]
		want int
	}{
		{"err err", result.Err[int]("string"), result.Err[int]("string"), 0},
		{"err less", result.Err[int]("a"), result.Err[int]("b"), -1},
		{"err greater", result.Err[int]("b"), result.Err[int]("a"), 1},
		{"zero err", result.Result[int, string]{}, result.Err[int]("string"), -1},
		{"zero empty err", result.Result[int, string]{}, result.Err[int](""), 0},
		{"err ok", result.Err[int]("string"), result.Ok[int, string](1), 1},
		{"ok err", result.Ok[int, string](1), result.Err[int]("string"), -1},
		{"ok less", result.Ok[int, string](1), result.Ok[int, string](2), -1},
		{"ok equal", result.Ok[int, string](2), result.Ok[int, string](2), 0},
		{"ok greater", result.Ok[int, string](3), result.Ok[int, string](1), 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := result.Compare(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b result.Result[int, string]
		want bool
	}{
		{"err err", result.Err[int]("string"), result.Err[int]("string"), true},
		{"err different", result.Err[int]("a"), result.Err[int]("b"), false},
		{"zero err", result.Result[int, string]{}, result.Err[int]("string"), false},
		{"zero empty err", result.Result[int, string]{}, result.Err[int](""), true},
		{"err ok", result.Err[int]("string"), result.Ok[int, string](0), false},
		{"ok err", result.Ok[int, string](0), result.Err[int]("string"), false},
		{"ok equal", result.Ok[int, string](2), result.Ok[int, string](2), true},
		{"ok different", result.Ok[int, string](2), result.Ok[int, string](3), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := result.Equal(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEqual_NaN(t *testing.T) {
	nan := result.Ok[float64, string](math.NaN())

	assert.False(t, result.Equal(nan, nan))
	assert.Equal(t, 0, result.Compare(nan, nan))
}

func TestGe(t *testing.T) {
	tests := []struct {
		name string
		a, b result.Result[int, string]
		want bool
	}{
		{"err err", result.Err[int]("string"), result.Err[int]("string"), true},
		{"err ok", result.Err[int]("string"), result.Ok[int, string](1), true},
		{"ok err", result.Ok[int, string](1), result.Err[int]("string"), false},
		{"ok less", result.Ok[int, string](1), result.Ok[int, string](2), false},
		{"ok equal", result.Ok[int, string](2), result.Ok[int, string](2), true},
		{"ok greater", result.Ok[int, string](3), result.Ok[int, string](1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := result.Ge(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGt(t *testing.T) {
	tests := []struct {
		name string
		a, b result.Result[int, string]
		want bool
	}{
		{"err err", result.Err[int]("string"), result.Err[int]("string"), false},
		{"err ok", result.Err[int]("string"), result.Ok[int, string](1), true},
		{"ok err", result.Ok[int, string](1), result.Err[int]("string"), false},
		{"ok less", result.Ok[int, string](1), result.Ok[int, string](2), false},
		{"ok equal", result.Ok[int, string](2), result.Ok[int, string](2), false},
		{"ok greater", result.Ok[int, string](3), result.Ok[int, string](1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := result.Gt(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLe(t *testing.T) {
	tests := []struct {
		name string
		a, b result.Result[int, string]
		want bool
	}{
		{"err err", result.Err[int]("string"), result.Err[int]("string"), true},
		{"err ok", result.Err[int]("string"), result.Ok[int, string](1), false},
		{"ok err", result.Ok[int, string](1), result.Err[int]("string"), true},
		{"ok less", result.Ok[int, string](1), result.Ok[int, string](2), true},
		{"ok equal", result.Ok[int, string](2), result.Ok[int, string](2), true},
		{"ok greater", result.Ok[int, string](3), result.Ok[int, string](1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := result.Le(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLt(t *testing.T) {
	tests := []struct {
		name string
		a, b result.Result[int, string]
		want bool
	}{
		{"err err", result.Err[int]("string"), result.Err[int]("string"), false},
		{"err ok", result.Err[int]("string"), result.Ok[int, string](1), false},
		{"ok err", result.Ok[int, string](1), result.Err[int]("string"), true},
		{"ok less", result.Ok[int, string](1), result.Ok[int, string](2), true},
		{"ok equal", result.Ok[int, string](2), result.Ok[int, string](2), false},
		{"ok greater", result.Ok[int, string](3), result.Ok[int, string](1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := result.Lt(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func ExampleCompare() {
	fmt.Println(result.Compare(result.Ok[int, string](1), result.Err[int]("e")))
	fmt.Println(result.Compare(result.Err[int]("e"), result.Ok[int, string](1)))
	fmt.Println(result.Compare(result.Ok[int, string](1), result.Ok[int, string](2)))
	fmt.Println(result.Compare(result.Ok[int, string](2), result.Ok[int, string](1)))
	fmt.Println(result.Compare(result.Ok[int, string](1), result.Ok[int, string](1)))
	fmt.Println(result.Compare(result.Err[int]("a"), result.Err[int]("a")))

	// Output:
	// -1
	// 1
	// -1
	// 1
	// 0
	// 0
}

func ExampleEqual() {
	fmt.Println(result.Equal(result.Ok[int, string](2), result.Ok[int, string](2)))
	fmt.Println(result.Equal(result.Ok[int, string](2), result.Err[int]("e")))
	fmt.Println(result.Equal(result.Err[int]("e"), result.Err[int]("e")))

	// Output:
	// true
	// false
	// true
}

func ExampleGe() {
	fmt.Println(result.Ge(result.Ok[int, string](2), result.Ok[int, string](1)))
	fmt.Println(result.Ge(result.Ok[int, string](1), result.Ok[int, string](1)))
	fmt.Println(result.Ge(result.Ok[int, string](1), result.Err[int]("e")))

	// Output:
	// true
	// true
	// false
}

func ExampleGt() {
	fmt.Println(result.Gt(result.Ok[int, string](2), result.Ok[int, string](1)))
	fmt.Println(result.Gt(result.Ok[int, string](1), result.Ok[int, string](1)))
	fmt.Println(result.Gt(result.Err[int]("e"), result.Ok[int, string](1)))

	// Output:
	// true
	// false
	// true
}

func ExampleLe() {
	fmt.Println(result.Le(result.Ok[int, string](1), result.Ok[int, string](2)))
	fmt.Println(result.Le(result.Ok[int, string](1), result.Ok[int, string](1)))
	fmt.Println(result.Le(result.Err[int]("e"), result.Ok[int, string](1)))

	// Output:
	// true
	// true
	// false
}

func ExampleLt() {
	fmt.Println(result.Lt(result.Ok[int, string](1), result.Ok[int, string](2)))
	fmt.Println(result.Lt(result.Ok[int, string](1), result.Ok[int, string](1)))
	fmt.Println(result.Lt(result.Ok[int, string](1), result.Err[int]("e")))

	// Output:
	// true
	// false
	// true
}
