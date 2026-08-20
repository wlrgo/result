package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlrgo/result"
)

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

func TestResult_Unpack(t *testing.T) {
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

func ExampleUnpack() {
	fmt.Println(result.Unpack(result.Ok[int, error](7)))
	fmt.Println(result.Unpack(result.Err[int](errors.New("late"))))

	// Output:
	// 7 <nil>
	// 0 late
}
