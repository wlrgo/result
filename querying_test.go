package result_test

import (
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
		pred func(int) bool
		want bool
	}{
		{"err", result.Err[int](ErrTest), true},
		{"ok", result.Ok[int, error](123), false},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.IsOk()
			assert.Equal(t, tt.want, got, "want: %v, got: %v", tt.want, got)
		})
	}
}
func TestResult_IsOkAnd(t *testing.T) {}
