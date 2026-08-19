package result_test

import (
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
		name  string
		give  result.Result[int, error]
		wantT int
		wantE error
	}{
		{"err", result.Err[int](ErrTest), 0, ErrTest},
		{"ok", result.Ok[int, error](7), 7, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, gotE := result.Unpack(tt.give)
			assert.Equal(t, tt.wantT, gotT)
			assert.Equal(t, tt.wantE, gotE)
		})
	}
}
