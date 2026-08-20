package result_test

import (
	"errors"
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
