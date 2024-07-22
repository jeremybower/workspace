package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathConverter(t *testing.T) {
	t.Parallel()

	pc := NewPathConverter("/source", "/target")
	targetPath, err := pc.SourceToTargetPath("/source/a")
	assert.NoError(t, err)
	assert.Equal(t, "/target/a", targetPath)

	targetPath, err = pc.SourceToTargetPath("/source/a/b")
	assert.NoError(t, err)
	assert.Equal(t, "/target/a/b", targetPath)

	_, err = pc.SourceToTargetPath("/invalid/a")
	assert.Error(t, err)

	targetPath, err = pc.TargetToSourcePath("/target/a")
	assert.NoError(t, err)
	assert.Equal(t, "/source/a", targetPath)

	targetPath, err = pc.TargetToSourcePath("/target/a/b")
	assert.NoError(t, err)
	assert.Equal(t, "/source/a/b", targetPath)

	_, err = pc.TargetToSourcePath("/invalid/a")
	assert.Error(t, err)
}
