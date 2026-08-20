package result_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlrgo/result"
)

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
