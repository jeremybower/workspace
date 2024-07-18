package funcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequire(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		value       any
		expectedErr error
	}{
		{
			name:        "nil",
			value:       nil,
			expectedErr: ErrValueRequired,
		},
		{
			name:        "Integer zero",
			value:       0,
			expectedErr: ErrValueRequired,
		},
		{
			name:        "Integer non-zero",
			value:       1,
			expectedErr: nil,
		},
		{
			name:        "String empty",
			value:       "",
			expectedErr: ErrValueRequired,
		},
		{
			name:        "String non-empty",
			value:       "a",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := Require()
			v, err := require(tt.value)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, v)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.value, v)
			}
		})
	}
}
