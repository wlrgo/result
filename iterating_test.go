package result_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlrgo/result"
)

func TestCollect(t *testing.T) {
	tests := []struct {
		name      string
		give      []result.Result[int, error]
		wantValue []int
		wantErr   bool
	}{
		{"empty", nil, []int{}, false},
		{
			"all ok",
			[]result.Result[int, error]{result.Ok[int, error](1), result.Ok[int, error](2)},
			[]int{1, 2},
			false,
		},
		{
			"first err",
			[]result.Result[int, error]{result.Err[int](errors.New("123")), result.Ok[int, error](2)},
			nil,
			true,
		},
		{
			"middle err",
			[]result.Result[int, error]{result.Ok[int, error](1), result.Err[int](errors.New("123")), result.Ok[int, error](3)},
			nil,
			true,
		},
		{
			"last err",
			[]result.Result[int, error]{result.Ok[int, error](1), result.Err[int](errors.New("123"))},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := result.Collect(tt.give)

			if tt.wantErr {
				assert.True(t, got.IsErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}

func TestResult_Seq(t *testing.T) {
	tests := []struct {
		name string
		give result.Result[int, error]
		want []int
	}{
		{"err", result.Err[int](errors.New("123")), nil},
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
}
