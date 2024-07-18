package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleTemplate(t *testing.T) {
	t.Parallel()

	// Create the template.
	templateFile, err := os.CreateTemp("", "template-*")
	require.NoError(t, err)
	defer templateFile.Close()

	_, err = templateFile.WriteString("Hello, {{ .Name }}!")
	require.NoError(t, err)

	// Create the config.
	config := map[string]interface{}{
		"Name": "World",
	}

	// Load the template.
	paths := []string{templateFile.Name()}
	tmpl, err := NewTemplate(paths)
	require.NoError(t, err)

	// Execute the template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, config)
	require.NoError(t, err)

	// Compare the output
	assert.Equal(t, "Hello, World!", buf.String())
}

func TestListTemplates(t *testing.T) {
	t.Parallel()

	// Create the included file.
	includeFile, err := os.CreateTemp("", "include-*")
	require.NoError(t, err)
	defer includeFile.Close()

	// Create the template.
	templateFile, err := os.CreateTemp("", "template-*")
	require.NoError(t, err)
	defer templateFile.Close()

	_, err = templateFile.WriteString(`{{ range listTemplates }}|{{ . }}|{{ end }}`)
	require.NoError(t, err)

	// Load the template.
	paths := []string{templateFile.Name(), includeFile.Name()}
	tmpl, err := NewTemplate(paths)
	require.NoError(t, err)

	// Execute the template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	require.NoError(t, err)

	// Compare the output
	expected := fmt.Sprintf(`|%s||%s|`, includeFile.Name(), templateFile.Name())
	assert.Equal(t, expected, buf.String())
}

func TestGlobFilter(t *testing.T) {
	t.Parallel()

	// Create the template.
	templateFile, err := os.CreateTemp("", "template-*")
	require.NoError(t, err)
	defer templateFile.Close()

	_, err = templateFile.WriteString(`{{ range (list "exclude" "include-1" "exclude" "include-2") | globFilter "include-*" }}|{{ . }}|{{ end }}`)
	require.NoError(t, err)

	// Load the template.
	paths := []string{templateFile.Name()}
	tmpl, err := NewTemplate(paths)
	require.NoError(t, err)

	// Execute the template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	require.NoError(t, err)

	// Compare the output
	assert.Equal(t, "|include-1||include-2|", buf.String())
}

func TestInclude(t *testing.T) {
	t.Parallel()

	// Create the included file.
	includeFile, err := os.CreateTemp("", "include-*")
	require.NoError(t, err)
	defer includeFile.Close()

	_, err = includeFile.WriteString("Hello, {{ .Name }}!")
	require.NoError(t, err)

	// Create the template.
	templateFile, err := os.CreateTemp("", "template-*")
	require.NoError(t, err)
	defer templateFile.Close()

	_, err = templateFile.WriteString(fmt.Sprintf(`{{ include "%s" . }}`, includeFile.Name()))
	require.NoError(t, err)

	// Create the config.
	config := map[string]interface{}{
		"Name": "World",
	}

	// Load the template.
	paths := []string{templateFile.Name(), includeFile.Name()}
	tmpl, err := NewTemplate(paths)
	require.NoError(t, err)

	// Execute the template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, config)
	require.NoError(t, err)

	// Compare the output
	assert.Equal(t, "Hello, World!", buf.String())
}

func TestAddFunc(t *testing.T) {
	t.Parallel()

	fm := template.FuncMap{}
	assert.NotPanics(t, func() {
		addFunc(fm, "test", func() {})
	})
	assert.Panics(t, func() {
		addFunc(fm, "test", func() {})
	})
}
