package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Data struct {
	Config map[string]any `yaml:"Config"`
}

// NewData reads the YAML file at the given path and returns a Data struct
func NewData(path string) (*Data, error) {
	// Read the file at the given path
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML data into a map
	var data Data
	err = yaml.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}

	// Ensure that the Config field is not nil
	if data.Config == nil {
		return nil, fmt.Errorf("required field not found: Config")
	}

	// Success
	return &data, nil
}

// Merge merges the given Data struct into this one
func (d *Data) Merge(other *Data) {
	if d.Config == nil {
		d.Config = make(map[string]any)
	}

	d.Config = mergeMaps(d.Config, other.Config)
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
