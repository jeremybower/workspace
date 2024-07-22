package internal

import (
	"fmt"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type ConfigSpecData struct {
	Config map[string]any `yaml:"Config"`
}

type ConfigSpec struct {
	fs     afero.Fs
	config map[string]any
}

func NewConfigSpec(fs afero.Fs, names []string) (*ConfigSpec, error) {
	configSpec := &ConfigSpec{
		fs:     fs,
		config: make(map[string]any),
	}

	for _, name := range names {
		err := configSpec.Merge(name)
		if err != nil {
			return nil, err
		}
	}

	return configSpec, nil
}

func (c *ConfigSpec) Merge(name string) error {
	// Read the file at the given path
	b, err := afero.ReadFile(c.fs, name)
	if err != nil {
		return err
	}

	// Unmarshal the YAML data into a map
	var data ConfigSpecData
	err = yaml.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	// Check that the required Config element is present
	if data.Config == nil {
		return fmt.Errorf("%w: required field '%s' not found", ErrConfigInvalid, "Config")
	}

	// Merge the maps.
	c.config = mergeMaps(c.config, data.Config)

	// Success
	return nil
}

func mergeMaps(a, b map[string]any) map[string]any {
	out := make(map[string]any, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[string]any); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[string]any); ok {
					out[k] = mergeMaps(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}
