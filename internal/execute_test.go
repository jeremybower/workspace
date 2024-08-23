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

func TestExecuteWhenSimilarPaths(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	th := NewTestHarness(t, fs)
	dir := th.TempDir()

	th.Mkdir(path.Join(dir, "modules"), 0755)
	th.Mkdir(path.Join(dir, "modules", "git"), 0755)
	th.WriteFileString(path.Join(dir, "modules", "core", "git", "module.dockerfile.tmpl"), "git-module\n")
	th.WriteFileString(path.Join(dir, "modules", "core", "git", "post-create.sh.tmpl"), "git-post-create\n")
	th.Mkdir(path.Join(dir, "modules", "git-lfs"), 0755)
	th.WriteFileString(path.Join(dir, "modules", "core", "git-lfs", "module.dockerfile.tmpl"), "git-lfs-module\n")
	th.WriteFileString(path.Join(dir, "modules", "core", "git-lfs", "post-create.sh.tmpl"), "git-lfs-post-create\n")

	th.WriteFileString(path.Join(dir, "Dockerfile.tmpl"),
		`{{ range $module := .Modules -}}{{ $path := printf "/modules/%s/module.dockerfile.tmpl" $module -}}{{ include $path $ }}{{ end -}}`)

	configFilename := path.Join(th.TempDir(), "config.yaml")
	th.WriteFileString(configFilename, "Config:\n  Modules:\n    - core/git\n    - core/git-lfs\n")

	outFilename := path.Join(th.TempDir(), "out")
	s, result := th.ExecuteString("/Dockerfile.tmpl", []string{
		path.Join(dir, "modules") + ":" + path.Join("/", "modules"),
		path.Join(dir, "Dockerfile.tmpl") + ":" + path.Join("/", "Dockerfile.tmpl"),
	}, []string{configFilename}, outFilename, DefaultOptions())

	assert.Equal(t, "git-module\ngit-lfs-module\n", s)
	assert.ElementsMatch(t, []string{outFilename}, result.Filenames)
	assert.GreaterOrEqual(t, result.Duration, time.Duration(0))
}
