package internal

import (
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigSpec(t *testing.T) {
	t.Parallel()

	// Prepare the test.
	th := NewTestHarness(t, afero.NewMemMapFs())
	dir := th.TempDir()

	// Write the first config file.
	config1 := path.Join(dir, "config1.yaml")
	th.WriteFileString(config1, `Config:
  Go:
    Version:
      Full: "1.22.5"`)

	// Write the second config file.
	config2 := path.Join(dir, "config2.yaml")
	th.WriteFileString(config2, `Config:
  Go:
    Packages:
      - something`)

	spec := th.NewConfigSpec(config1, config2)
	assert.Equal(t, map[string]any{
		"Go": map[string]any{
			"Version": map[string]any{
				"Full": "1.22.5",
			},
			"Packages": []any{
				"something",
			},
		},
	}, spec.config)
}

func TestMergeMaps(t *testing.T) {
	t.Parallel()

	result := mergeMaps(
		map[string]any{
			"Config": map[string]any{
				"Go": map[string]any{
					"Version": map[string]any{
						"Full": "1.22.5",
					},
				},
			},
		},
		map[string]any{
			"Config": map[string]any{
				"Go": map[string]any{
					"Packages": []any{
						"something",
					},
				},
			},
		},
	)

	assert.Equal(t, map[string]any{
		"Config": map[string]any{
			"Go": map[string]any{
				"Version": map[string]any{
					"Full": "1.22.5",
				},
				"Packages": []any{
					"something",
				},
			},
		},
	}, result)
}
