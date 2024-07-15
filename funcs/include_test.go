package funcs

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInclude(t *testing.T) {
	t.Parallel()

	tmpl1, err := template.New("tmpl1").Parse("tmpl1: {{ . }}")
	require.NoError(t, err)

	tmpl2, err := template.New("tmpl2").Parse("tmpl2: {{ . }}")
	require.NoError(t, err)

	f := Include(map[string]*template.Template{
		"tmpl1": tmpl1,
		"tmpl2": tmpl2,
	})

	result, err := f("tmpl1", "foo")
	require.NoError(t, err)
	assert.Equal(t, "tmpl1: foo", result)

	result, err = f("tmpl2", "bar")
	require.NoError(t, err)
	assert.Equal(t, "tmpl2: bar", result)
}
