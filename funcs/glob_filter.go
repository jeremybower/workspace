package funcs

import (
	"fmt"
	"path/filepath"
)

var ErrUnsupportedType = fmt.Errorf("unsupported type")

func GlobFilter() func(string, []any) ([]any, error) {
	return func(pattern string, items []any) ([]any, error) {
		filtered := []any{}
		for _, item := range items {
			switch v := item.(type) {
			case string:
				if match, err := filepath.Match(pattern, v); err != nil {
					return nil, err
				} else if match {
					filtered = append(filtered, v)
				}
			default:
				return nil, fmt.Errorf("%w: %T", ErrUnsupportedType, item)
			}
		}

		return filtered, nil
	}
}
