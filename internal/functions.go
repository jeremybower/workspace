package internal

import (
	"fmt"
	"path"
	"text/template"
)

var DummyFunctions = &Functions{}

type Functions struct {
	filename string
	mounts   Mounts
	cache    *TemplateCache
}

func NewFunctions(
	filename string,
	mounts Mounts,
	cache *TemplateCache,
) *Functions {
	return &Functions{
		filename: filename,
		mounts:   mounts,
		cache:    cache,
	}
}

func (f *Functions) FuncMap() template.FuncMap {
	return template.FuncMap{
		"dirs":        f.dirsFunc,
		"filename":    f.filenameFunc,
		"files":       f.filesFunc,
		"include":     f.includeFunc,
		"includeText": f.includeTextFunc,
	}
}

func (f *Functions) dirsFunc(pattern string) ([]string, error) {
	return f.mounts.Directories(pattern)
}

func (f *Functions) filenameFunc() string {
	return f.filename
}

func (f *Functions) filesFunc(pattern string) ([]string, error) {
	return f.mounts.Files(pattern)
}

func (f *Functions) includeFunc(filename string, data any) (string, error) {
	// Check if the filename is a relative path.
	if !path.IsAbs(filename) {
		// Convert to an absolute path using the directory of
		// the current filename as the root.
		filename = path.Join(path.Dir(f.filename), filename)

		// Clean to remove any ".." or ".".
		filename = path.Clean(filename)
	}

	// Load the template from the cache.
	t, err := f.cache.Template(filename)
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, filename)
	}

	// Execute the template.
	return t.ExecuteString(f.mounts, data)
}

func (f *Functions) includeTextFunc(filename string) (string, error) {
	// Check if the filename is a relative path.
	if !path.IsAbs(filename) {
		// Convert to an absolute path using the directory of
		// the current filename as the root.
		filename = path.Join(path.Dir(f.filename), filename)

		// Clean to remove any ".." or ".".
		filename = path.Clean(filename)
	}

	// Read the file as a string.
	return f.mounts.ReadFileString(filename)
}
