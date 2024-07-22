package internal

import (
	"fmt"
	"path"
	"sync"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type TemplateCache struct {
	mounts    Mounts
	templates map[string]*template.Template
	mutex     sync.RWMutex
}

func NewTemplateCache(mounts Mounts) *TemplateCache {
	return &TemplateCache{
		mounts:    mounts,
		templates: make(map[string]*template.Template),
	}
}

func (cache *TemplateCache) Template(name string) (*Template, error) {
	// The template name must be an absolute path.
	if !path.IsAbs(name) {
		return nil, fmt.Errorf("%w: template: %s", ErrAbsolutePathRequired, name)
	}

	// Define a function to safely readCached a template from the cache.
	readCached := func() *template.Template {
		// Lock the cache for reading.
		cache.mutex.RLock()
		defer cache.mutex.RUnlock()

		// Return the template if it is cached, nil otherwise.
		return cache.templates[name]
	}

	// Define a function to safely createAndCache and cache the template.
	createAndCache := func() (*template.Template, error) {
		// Acquire a write lock.
		cache.mutex.Lock()
		defer cache.mutex.Unlock()

		// Check if the template is already loaded.
		t := cache.templates[name]
		if t == nil {
			// Read the template from the mounts.
			s, err := cache.mounts.ReadFileString(name)
			if err != nil {
				return nil, err
			}

			// Create and parse the template.
			t, err = template.New(name).Funcs(sprig.FuncMap()).Funcs(DummyFunctions.FuncMap()).Parse(s)
			if err != nil {
				return nil, err
			}

			// Save the template for reuse.
			cache.templates[name] = t
		}

		// Return the template.
		return t, nil
	}

	// Get the cached template.
	t := readCached()
	if t == nil {
		var err error
		t, err = createAndCache()
		if err != nil {
			return nil, err
		}
	}

	// Clone the template.
	t, err := t.Clone()
	if err != nil {
		return nil, err
	}

	// Success.
	return NewTemplate(t, cache), nil
}
