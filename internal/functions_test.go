package internal

import (
	"fmt"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirsFunc(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	th := NewTestHarness(t, fs)
	dir := th.TempDir()

	th.WriteFileString(path.Join(dir, "a"), "")
	th.WriteFileString(path.Join(dir, "b"), "")
	th.Mkdir(path.Join(dir, "c"), 0755)
	th.WriteFileString(path.Join(dir, "c", "d"), "")
	th.WriteFileString(path.Join(dir, "c", "e"), "")
	th.WriteFileString(path.Join(dir, "f"), "")

	source := dir
	target := "/target"
	spec := source + ":" + target

	filename := "/target/filename"
	mounts := th.NewMounts(spec)
	cache := NewTemplateCache(mounts, DefaultOptions())
	funcs := NewFunctions(filename, mounts, cache)

	dirs, err := funcs.dirsFunc("/target/*")
	require.NoError(t, err)
	assert.ElementsMatch(t, dirs, []string{
		"/target/c",
	})
}

func TestFilenameFunc(t *testing.T) {
	t.Parallel()

	filename := "/target/filename"
	funcs := NewFunctions(filename, nil, nil)

	assert.Equal(t, filename, funcs.filenameFunc())
}

func TestFilesFunc(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	th := NewTestHarness(t, fs)
	dir := th.TempDir()

	th.WriteFileString(path.Join(dir, "a"), "")
	th.WriteFileString(path.Join(dir, "b"), "")
	th.Mkdir(path.Join(dir, "c"), 0755)
	th.WriteFileString(path.Join(dir, "c", "d"), "")
	th.WriteFileString(path.Join(dir, "c", "e"), "")
	th.WriteFileString(path.Join(dir, "f"), "")

	source := dir
	target := "/target"
	spec := source + ":" + target

	filename := "/target/filename"
	mounts := th.NewMounts(spec)
	cache := NewTemplateCache(mounts, DefaultOptions())
	funcs := NewFunctions(filename, mounts, cache)

	files, err := funcs.filesFunc("/target/*")
	require.NoError(t, err)
	assert.ElementsMatch(t, files, []string{
		"/target/a",
		"/target/b",
		"/target/f",
	})

	files, err = funcs.filesFunc("/target/c/*")
	require.NoError(t, err)
	assert.ElementsMatch(t, files, []string{
		"/target/c/d",
		"/target/c/e",
	})
}

func TestIncludeFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		includePath    string
		includeContent string
		data           any
		expected       string
	}{
		{
			name:           "include absolue template",
			includePath:    "/target/b",
			includeContent: `{{ .Name }}`,
			data:           map[string]interface{}{"Name": "absolute"},
			expected:       "Hello, absolute!",
		},
		{
			name:           "include relative template",
			includePath:    "./b",
			includeContent: `{{ .Name }}`,
			data:           map[string]interface{}{"Name": "relative"},
			expected:       "Hello, relative!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			th := NewTestHarness(t, fs)
			dir := th.TempDir()

			contentA := fmt.Sprintf(`Hello, {{ include "%s" . }}!`, tt.includePath)

			th.WriteFileString(path.Join(dir, "a"), contentA)
			th.WriteFileString(path.Join(dir, "b"), tt.includeContent)

			source := dir
			target := "/target"
			spec := source + ":" + target

			filename := "/target/filename"
			mounts := th.NewMounts(spec)
			cache := NewTemplateCache(mounts, DefaultOptions())
			funcs := NewFunctions(filename, mounts, cache)

			content, err := funcs.includeFunc("/target/a", tt.data)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, content)
		})
	}
}

func TestIncludeTextFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		includePath    string
		includeContent string
		expected       string
	}{
		{
			name:           "include absolue template",
			includePath:    "/target/b",
			includeContent: "absolute",
			expected:       "Hello, absolute!",
		},
		{
			name:           "include relative template",
			includePath:    "./b",
			includeContent: "relative",
			expected:       "Hello, relative!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			th := NewTestHarness(t, fs)
			dir := th.TempDir()

			contentA := fmt.Sprintf(`Hello, {{ includeText "%s" }}!`, tt.includePath)

			th.WriteFileString(path.Join(dir, "a"), contentA)
			th.WriteFileString(path.Join(dir, "b"), tt.includeContent)

			source := dir
			target := "/target"
			spec := source + ":" + target

			filename := "/target/filename"
			mounts := th.NewMounts(spec)
			cache := NewTemplateCache(mounts, DefaultOptions())
			funcs := NewFunctions(filename, mounts, cache)

			content, err := funcs.includeFunc("/target/a", nil)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, content)
		})
	}
}
