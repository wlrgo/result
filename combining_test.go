package result_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlrgo/result"
)

func TestAnd(t *testing.T) {
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
			got := result.And(tt.a, tt.b)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, got.UnwrapErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.UnwrapOr("<none>"))
		})
	}
}

func TestAndThen(t *testing.T) {
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
			got := result.AndThen(tt.a, func(i int) result.Result[string, error] {
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

func TestOr(t *testing.T) {
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
			got := result.Or(tt.a, tt.b)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, got.UnwrapErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.UnwrapOr(-1))
		})
	}
}

func TestOrElse(t *testing.T) {
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
			got := result.OrElse(tt.give, func(err error) result.Result[int, error] {
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
