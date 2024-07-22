package internal

import (
	"io"
	"strings"
	"text/template"
)

type Template struct {
	t     *template.Template
	cache *TemplateCache
}

func NewTemplate(
	t *template.Template,
	cache *TemplateCache,
) *Template {
	return &Template{t, cache}
}

func (t *Template) Execute(wr io.Writer, mounts Mounts, data any) error {
	cloned, err := t.t.Clone()
	if err != nil {
		return err
	}

	filename := t.t.Name()
	funcs := NewFunctions(filename, mounts, t.cache)
	cloned.Funcs(funcs.FuncMap())

	return cloned.Execute(wr, data)
}

func (t *Template) ExecuteString(mounts Mounts, data any) (string, error) {
	buf := new(strings.Builder)
	err := t.Execute(buf, mounts, data)
	return buf.String(), err
}
