package internal

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMountWhenDirectory(t *testing.T) {
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
	m := th.NewMount(spec)

	tests := []struct {
		pattern  string
		expected []string
	}{
		{
			pattern: "/target/a",
			expected: []string{
				"/target/a",
			},
		},
		{
			pattern: "/target/*",
			expected: []string{
				"/target/a",
				"/target/b",
				"/target/f",
			},
		},
		{
			pattern: "/target/c/*",
			expected: []string{
				"/target/c/d",
				"/target/c/e",
			},
		},
		{
			pattern:  "/target/z/*",
			expected: []string{},
		},
		{
			pattern:  "/invalid/*",
			expected: []string{},
		},
	}

	for _, test := range tests {
		var files []string
		err := m.Files(test.pattern, &files, nil)
		assert.NoError(t, err)
		assert.ElementsMatch(t, test.expected, files)
	}
}

func TestMountWhenFile(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	th := NewTestHarness(t, fs)
	dir := th.TempDir()
	th.WriteFileString(path.Join(dir, "a"), "")

	source := path.Join(dir, "a")
	target := "/target/x"
	spec := source + ":" + target
	m := th.NewMount(spec)

	tests := []struct {
		pattern  string
		expected []string
	}{
		{
			pattern: "/target/x",
			expected: []string{
				"/target/x",
			},
		},
		{
			pattern: "/target/*",
			expected: []string{
				"/target/x",
			},
		},
		{
			pattern:  "/target/z",
			expected: []string{},
		},
		{
			pattern:  "/invalid/*",
			expected: []string{},
		},
	}

	for _, test := range tests {
		var files []string
		err := m.Files(test.pattern, &files, nil)
		assert.NoError(t, err)
		assert.ElementsMatch(t, test.expected, files)
	}
}

func TestMountWhenFileMountedAsDirectory(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	th := NewTestHarness(t, fs)
	dir := th.TempDir()
	th.WriteFileString(path.Join(dir, "a"), "")

	source := appendPathSeparator(path.Join(dir, "a"))
	target := "/target"
	spec := source + ":" + target
	err := th.NewMountExpectingError(spec)
	require.ErrorIs(t, err, ErrMountInvalid)
}

func TestMountWhenInvalid(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	specs := []string{"source", "source:", ":target", "source:target:foo"}
	for _, spec := range specs {
		t.Run(spec, func(t *testing.T) {
			_, err := NewMount(fs, spec)
			require.ErrorIs(t, err, ErrMountInvalid)
		})
	}
}

func TestMountWhenSourceNotFound(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	th := NewTestHarness(t, fs)
	dir := th.TempDir()

	source := path.Join(dir, "a")
	target := "/target"
	spec := source + ":" + target

	err := th.NewMountExpectingError(spec)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestMountWhenTargetNotAbsolute(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	th := NewTestHarness(t, fs)
	dir := th.TempDir()
	th.WriteFileString(path.Join(dir, "a"), "")

	source := path.Join(dir, "a")
	target := path.Join("..", "target")
	spec := source + ":" + target

	err := th.NewMountExpectingError(spec)
	require.ErrorIs(t, err, ErrAbsolutePathRequired)
}

func TestMountFileAsChildOfTarget(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	th := NewTestHarness(t, fs)
	dir := th.TempDir()
	th.WriteFileString(path.Join(dir, "a"), "")

	source := path.Join(dir, "a")
	target := "/target/"
	spec := source + ":" + target
	m := th.NewMount(spec)

	var files []string
	err := m.Files("/target/*", &files, nil)
	assert.NoError(t, err)
	assert.ElementsMatch(t, files, []string{"/target/a"})
}

func TestListFiles(t *testing.T) {
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

	var files []string
	err := listFiles(fs, dir, &files, func(p string) (string, error) {
		return p, nil
	})
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{
		path.Join(dir, "a"),
		path.Join(dir, "b"),
		path.Join(dir, "c", "d"),
		path.Join(dir, "c", "e"),
		path.Join(dir, "f"),
	}, files)
}
