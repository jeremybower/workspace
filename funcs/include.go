package funcs

import (
	"fmt"
	"strings"
	"text/template"
)

func Include(templates map[string]*template.Template) func(string, interface{}) (string, error) {
	return func(includedPath string, data interface{}) (string, error) {
		t, ok := templates[includedPath]
		if !ok {
			return "", fmt.Errorf("template not found: %s", includedPath)
		}

		var buf strings.Builder
		err := t.Execute(&buf, data)
		if err != nil {
			return "", err
		}

		return buf.String(), nil
	}
}
