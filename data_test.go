package tmpl

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestData(t *testing.T) {
	t.Parallel()

	// Data 1
	file1, err := os.CreateTemp("", "data-*")
	require.NoError(t, err)

	_, err = file1.WriteString(`Config:
  Go:
    Version:
      Full: "1.22.5"`)
	require.NoError(t, err)

	data1, err := NewData(file1.Name())
	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"Go": map[string]any{
			"Version": map[string]any{
				"Full": "1.22.5",
			},
		},
	}, data1.Config)

	// Data 2
	file2, err := os.CreateTemp("", "data-*")
	require.NoError(t, err)

	_, err = file2.WriteString(`Config:
  Go:
    Packages:
      - something`)
	require.NoError(t, err)

	data2, err := NewData(file2.Name())
	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"Go": map[string]any{
			"Packages": []any{
				"something",
			},
		},
	}, data2.Config)

	data1.Merge(data2)
	assert.Equal(t, map[string]any{
		"Go": map[string]any{
			"Version": map[string]any{
				"Full": "1.22.5",
			},
			"Packages": []any{
				"something",
			},
		},
	}, data1.Config)
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
