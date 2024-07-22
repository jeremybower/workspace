package internal

import (
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestMountsFiles(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	th := NewTestHarness(t, fs)

	dir1 := th.TempDir()
	th.WriteFileString(path.Join(dir1, "a"), "")
	th.Mkdir(path.Join(dir1, "b"), 0755)
	th.WriteFileString(path.Join(dir1, "b", "c"), "")
	th.Mkdir(path.Join(dir1, "0"), 0755)
	th.WriteFileString(path.Join(dir1, "0", "9"), "")

	dir2 := th.TempDir()
	th.WriteFileString(path.Join(dir2, "d"), "")

	dir3 := th.TempDir()
	th.WriteFileString(path.Join(dir3, "1"), "")
	th.WriteFileString(path.Join(dir3, "2"), "")

	spec1 := dir1 + ":/target"
	spec2 := path.Join(dir2, "d") + ":/target/d"
	spec3 := dir3 + ":/target/0"

	mounts := th.NewMounts(spec1, spec2, spec3)
	files, err := mounts.Files("/target/*")
	assert.NoError(t, err)
	assert.ElementsMatch(t, files, []string{
		"/target/a",
		"/target/d",
	})

	files, err = mounts.Files("/target/b/*")
	assert.NoError(t, err)
	assert.ElementsMatch(t, files, []string{
		"/target/b/c",
	})

	files, err = mounts.Files("/target/0/*")
	assert.NoError(t, err)
	assert.ElementsMatch(t, files, []string{
		"/target/0/1",
		"/target/0/2",
	})
}
