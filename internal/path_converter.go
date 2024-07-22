package internal

import (
	"errors"
	"path/filepath"
	"strings"
)

type PathConverter struct {
	sourcePath string
	targetPath string
}

func NewPathConverter(sourcePath, targetPath string) *PathConverter {
	return &PathConverter{
		sourcePath: sourcePath,
		targetPath: targetPath,
	}
}

func (p *PathConverter) SourceToTargetPath(sourcePath string) (string, error) {
	if !strings.HasPrefix(sourcePath, p.sourcePath) {
		return "", errors.New("source path is not in the source directory")
	}
	return filepath.Join(p.targetPath, sourcePath[len(p.sourcePath):]), nil
}

func (p *PathConverter) TargetToSourcePath(targetPath string) (string, error) {
	if !strings.HasPrefix(targetPath, p.targetPath) {
		return "", errors.New("target path is not in the target directory")
	}
	return filepath.Join(p.sourcePath, targetPath[len(p.targetPath):]), nil
}
