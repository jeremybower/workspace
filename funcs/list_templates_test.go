package funcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTemp(t *testing.T) {
	t.Parallel()

	f := ListTemplates([]string{"a", "b", "c"})
	result := f()
	assert.Equal(t, []any{"a", "b", "c"}, result)
}
