package internal

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

type TestHarness struct {
	t  *testing.T
	fs afero.Fs
}

func NewTestHarness(t *testing.T, fs afero.Fs) *TestHarness {
	return &TestHarness{t, fs}
}

//-----------------------------------------------------------------------------
// Config specs
//-----------------------------------------------------------------------------

func (th *TestHarness) NewConfigSpec(names ...string) *ConfigSpec {
	configSpec, err := NewConfigSpec(th.fs, names)
	require.NoError(th.t, err)
	require.NotNil(th.t, configSpec)
	return configSpec
}

//-----------------------------------------------------------------------------
// Filesystem
//-----------------------------------------------------------------------------

func (th *TestHarness) Mkdir(name string, perm os.FileMode) {
	err := th.fs.Mkdir(name, perm)
	require.NoError(th.t, err)
}

func (th *TestHarness) ReadFileString(name string) string {
	b, err := afero.ReadFile(th.fs, name)
	require.NoError(th.t, err)
	return string(b)
}

func (th *TestHarness) TempDir() string {
	dir, err := afero.TempDir(th.fs, "", "")
	require.NoError(th.t, err)
	require.NotEmpty(th.t, dir)
	return dir
}

func (th *TestHarness) WriteFileString(name, content string) {
	err := afero.WriteFile(th.fs, name, []byte(content), 0644)
	require.NoError(th.t, err)
}

//-----------------------------------------------------------------------------
// Mounts
//-----------------------------------------------------------------------------

func (th *TestHarness) NewMount(spec string) *Mount {
	mount, err := NewMount(th.fs, spec)
	require.NoError(th.t, err)
	require.NotNil(th.t, mount)
	return mount
}

func (th *TestHarness) NewMountExpectingError(spec string) error {
	mount, err := NewMount(th.fs, spec)
	require.Error(th.t, err)
	require.Nil(th.t, mount)
	return err
}

func (th *TestHarness) NewMounts(specs ...string) Mounts {
	mounts, err := NewMounts(th.fs, specs)
	require.NoError(th.t, err)
	require.NotNil(th.t, mounts)
	return mounts
}

//-----------------------------------------------------------------------------
// Templates
//-----------------------------------------------------------------------------

func (th *TestHarness) Execute(tmplFilename string, mountSpecs []string, configFilenames []string, outFilename string, opts Options) *Result {
	result, err := Execute(th.fs, tmplFilename, mountSpecs, configFilenames, outFilename, opts)
	require.NoError(th.t, err)
	require.NotNil(th.t, result)
	return result
}

func (th *TestHarness) ExecuteString(tmplFilename string, mountSpecs []string, configFilenames []string, outFilename string, opts Options) (string, *Result) {
	result := th.Execute(tmplFilename, mountSpecs, configFilenames, outFilename, opts)
	return th.ReadFileString(outFilename), result
}
