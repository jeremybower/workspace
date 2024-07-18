package internal

import (
	"io"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/jeremybower/tmpl/internal/funcs"
)

type Template struct {
	t *template.Template
}

func NewTemplate(
	paths []string,
) (*Template, error) {
	// Create a uninitialized template functions to parse the templates.
	fm := templateFuncs(nil, nil)

	// Parse the included templates.
	templates := make(map[string]*template.Template)
	for _, path := range paths {
		t, err := parseTemplate(path, fm)
		if err != nil {
			return nil, err
		}

		templates[path] = t
	}

	// Create initialized template functions.
	fm = templateFuncs(paths, templates)

	// Update the functions on the included templates.
	for _, t := range templates {
		t.Funcs(fm)
	}

	// Success.
	return &Template{templates[paths[0]]}, nil
}

func (t *Template) Execute(wr io.Writer, data interface{}) error {
	return t.t.Execute(wr, data)
}

func parseTemplate(path string, fm template.FuncMap) (*template.Template, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return template.New(path).Funcs(fm).Parse(string(b))
}

func templateFuncs(
	paths []string,
	templates map[string]*template.Template,
) template.FuncMap {
	fm := sprig.FuncMap()

	addFunc(fm, "globFilter", funcs.GlobFilter())
	addFunc(fm, "include", funcs.Include(templates))
	addFunc(fm, "listTemplates", funcs.ListTemplates(paths))
	addFunc(fm, "require", funcs.Require())

	return fm
}

func addFunc(fm template.FuncMap, name string, f interface{}) {
	if fm[name] != nil {
		panic("function already exists")
	}

	fm[name] = f
}
