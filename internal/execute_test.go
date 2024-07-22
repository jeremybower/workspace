package internal

import (
	"path"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	th := NewTestHarness(t, fs)
	dir := th.TempDir()

	th.WriteFileString(path.Join(dir, "a"), `{{ include "./b/c" . }}`)
	th.Mkdir(path.Join(dir, "b"), 0755)
	th.WriteFileString(path.Join(dir, "b", "c"), "Data: {{ .Test }}\nFilename: {{ filename }}\nDirs: {{ dirs \"/target/*\" }}\nFiles: {{ files \"/target/*\" }}")

	source := dir
	target := path.Join("/", "target")
	spec := source + ":" + target

	configFilename := path.Join(th.TempDir(), "config.yaml")
	th.WriteFileString(configFilename, "Config:\n  Test: Hello, World!")

	outFilename := path.Join(th.TempDir(), "out")
	s, result := th.ExecuteString("/target/a", []string{spec}, []string{configFilename}, outFilename, DefaultOptions())

	assert.Equal(t, "Data: Hello, World!\nFilename: /target/b/c\nDirs: [/target/b]\nFiles: [/target/a]", s)
	assert.ElementsMatch(t, []string{outFilename}, result.Filenames)
	assert.GreaterOrEqual(t, result.Duration, time.Duration(0))
}
