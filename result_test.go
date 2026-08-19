package result_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wlrgo/result"
)

var ErrTest = errors.New("test")

func TestFrom(t *testing.T) {
	tests := []struct {
		name      string
		v         int
		err       error
		wantValue int
		wantErr   bool
	}{
		{"nil error", 7, nil, 7, false},
		{"error", 7, ErrTest, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := result.From(tt.v, tt.err)

			if tt.wantErr {
				assert.True(t, got.IsErr())
				return
			}

			assert.Equal(t, tt.wantValue, got.Unwrap())
		})
	}
}
