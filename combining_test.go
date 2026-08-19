package result_test

import (
	"errors"
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
		wantErr   bool
	}{
		{
			"ok err",
			result.Ok[int, error](2),
			result.Err[string](errors.New("123")),
			"",
			true,
		},
		{
			"err ok",
			result.Err[int](errors.New("123")),
			result.Ok[string, error]("2"),
			"",
			true,
		},
		{
			"ok ok",
			result.Ok[int, error](2),
			result.Ok[string, error]("2"),
			"2",
			false,
		},
		{
			"err err",
			result.Err[int](errors.New("123")),
			result.Err[string](errors.New("str")),
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := result.And(tt.a, tt.b)

			if tt.wantErr {
				assert.True(t, got.IsErr())
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
		wantValue string
		wantErr   bool
		wantCalls int
	}{
		{
			"ok",
			result.Ok[int, error](2),
			"4",
			false,
			1,
		},
		{
			"err",
			result.Err[int](errors.New("123")),
			"",
			true,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := result.AndThen(tt.a, func(i int) result.Result[string, error] {
				calls++
				return result.Ok[string, error](strconv.Itoa(i * i))
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantErr {
				assert.True(t, got.IsErr())
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
		wantErr   bool
	}{
		{
			"ok err",
			result.Ok[int, error](2),
			result.Err[int](errors.New("123")),
			2,
			false,
		},
		{
			"err ok",
			result.Err[int](errors.New("123")),
			result.Ok[int, error](67),
			67,
			false,
		},
		{
			"ok ok",
			result.Ok[int, error](2),
			result.Ok[int, error](67),
			2,
			false,
		},
		{
			"err err",
			result.Err[int](errors.New("123")),
			result.Err[int](errors.New("str")),
			-1,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := result.Or(tt.a, tt.b)

			if tt.wantErr {
				assert.True(t, got.IsErr())
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
		wantValue int
		wantErr   bool
		wantCalls int
	}{
		{
			"ok",
			result.Ok[int, error](2),
			2,
			false,
			0,
		},
		{
			"err",
			result.Err[int](errors.New("123")),
			67,
			false,
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got := result.OrElse(tt.give, func(error) result.Result[int, error] {
				calls++
				return result.Ok[int, error](67)
			})

			assert.Equal(t, tt.wantCalls, calls)

			if tt.wantErr {
				assert.True(t, got.IsErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.UnwrapOr(-1))
		})
	}
}
