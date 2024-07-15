package funcs

import (
	"slices"
)

func ListTemplates(paths []string) func() []any {
	paths = slices.Clone(paths)
	slices.Sort(paths)
	result := make([]any, len(paths))
	for i, path := range paths {
		result[i] = path
	}
	return func() []any {
		return result
	}
}
