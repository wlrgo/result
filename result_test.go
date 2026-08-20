package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlrgo/result"
)

var (
	ErrTest  = errors.New("test")
	ErrOther = errors.New("other")
)

type typedNilError struct{}

func (*typedNilError) Error() string { return "typed nil" }

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

func ExampleResult() {
	var res result.Result[int, error]

	fmt.Println(res.IsErr())
	// Output: true
}

func ExampleErr() {
	res := result.Err[int](errors.New("late"))

	fmt.Println(res.IsErr())
	// Output: true
}

func ExampleOk() {
	res := result.Ok[int, error](10)

	fmt.Println(res.IsOk())
	fmt.Println(res.Unwrap())
	// Output:
	// true
	// 10
}
