package funcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobFilter(t *testing.T) {
	t.Parallel()

	f := GlobFilter()
	result, err := f("*.txt", []any{"a.txt", "b.txt", "123", "c.txt", "456"})
	assert.NoError(t, err)
	assert.Equal(t, []any{"a.txt", "b.txt", "c.txt"}, result)
}

func TestGlobFilterUnsupportedTypes(t *testing.T) {
	t.Parallel()

	f := GlobFilter()
	result, err := f("*.txt", []any{123, 456})
	assert.ErrorIs(t, err, ErrUnsupportedType)
	assert.Nil(t, result)
}
