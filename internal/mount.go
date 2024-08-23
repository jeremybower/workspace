package internal

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/afero"
)

func NewMount(fs afero.Fs, spec string) (*Mount, error) {
	// Check the mount format.
	paths := strings.Split(spec, ":")
	if len(paths) != 2 || paths[0] == "" || paths[1] == "" {
		return nil, fmt.Errorf("%w: format must be 'source:target'", ErrMountInvalid)
	}

	// Check that the source path is valid.
	sourceEndsWithSeparator := strings.HasSuffix(paths[0], string(filepath.Separator))
	sourcePath, err := filepath.Abs(paths[0])
	if err != nil {
		return nil, err
	}

	// Check if the source exists.
	source, err := fs.Stat(sourcePath)
	if err != nil {
		return nil, err
	}

	// Check that the target path is valid.
	targetEndsWithSeparator := strings.HasSuffix(paths[1], string(filepath.Separator))
	targetPath := filepath.Clean(paths[1])
	if !path.IsAbs(targetPath) {
		return nil, fmt.Errorf("%w: mount target: %s", ErrAbsolutePathRequired, targetPath)
	}

	// Check if source file is mounted as a directory.
	if !source.IsDir() && sourceEndsWithSeparator {
		return nil, fmt.Errorf("%w: source file mounted as directory: %s", ErrMountInvalid, sourcePath)
	}

	// If the target path ends in a separator and the source path does not, then
	// add the source base name to the target so that it's mounted as a child of
	// the target.
	if targetEndsWithSeparator && !sourceEndsWithSeparator {
		targetPath = filepath.Join(targetPath, filepath.Base(sourcePath))
	}

	// Create a path converter.
	pathConverter := NewPathConverter(sourcePath, targetPath)

	// Check if file or directory.
	var targetDirs []string
	var targetFiles []string
	if source.IsDir() {
		// List directories in the mount.
		if err := listDirs(fs, sourcePath, &targetDirs, pathConverter.SourceToTargetPath); err != nil {
			return nil, err
		}

		// List files in the mount.
		if err := listFiles(fs, sourcePath, &targetFiles, pathConverter.SourceToTargetPath); err != nil {
			return nil, err
		}
	} else {
		targetFile, err := pathConverter.SourceToTargetPath(sourcePath)
		if err != nil {
			return nil, err
		}

		targetFiles = append(targetFiles, targetFile)
	}

	// Sort to enable binary search.
	slices.Sort(targetDirs)
	slices.Sort(targetFiles)

	// Create the mount.
	return &Mount{
		fs:            fs,
		sourcePath:    sourcePath,
		targetPath:    targetPath,
		targetDirs:    targetDirs,
		targetFiles:   targetFiles,
		pathConverter: pathConverter,
		directory:     source.IsDir(),
	}, nil
}

type Mount struct {
	fs            afero.Fs
	sourcePath    string
	targetPath    string
	targetDirs    []string
	targetFiles   []string
	pathConverter *PathConverter
	directory     bool
}

func (m *Mount) Dirs(pattern string, directories *[]string, excludeFns []func(string) bool) error {
	// Match files.
	for _, targetDir := range m.targetDirs {
		// Check if the target directory is excluded.
		if exclude(targetDir, excludeFns) {
			continue
		}

		// Check if the directory matches the pattern.
		if match, err := filepath.Match(pattern, targetDir); err != nil {
			return err
		} else if match {
			*directories = append(*directories, targetDir)
		}
	}

	// Success.
	return nil
}

func (m *Mount) Files(pattern string, files *[]string, excludeFns []func(string) bool) error {
	// Match files.
	for _, targetFile := range m.targetFiles {
		// Check if the target file is excluded.
		if exclude(targetFile, excludeFns) {
			continue
		}

		// Check if the file matches the pattern.
		if match, err := filepath.Match(pattern, targetFile); err != nil {
			return err
		} else if match {
			*files = append(*files, targetFile)
		}
	}

	// Success.
	return nil
}

func (m *Mount) ReadFileString(targetPath string) (string, error) {
	_, ok := slices.BinarySearch(m.targetFiles, targetPath)
	if !ok {
		return "", fmt.Errorf("%w: %s", os.ErrNotExist, targetPath)
	}

	sourcePath, err := m.pathConverter.TargetToSourcePath(targetPath)
	if err != nil {
		return "", err
	}

	b, err := afero.ReadFile(m.fs, sourcePath)
	if err != nil {
		return "", fmt.Errorf("error reading source file: %w: %s", err, sourcePath)
	}

	return string(b), nil
}

func listDirs(fs afero.Fs, sourcePath string, directories *[]string, pathConverter func(string) (string, error)) error {
	return afero.Walk(fs, sourcePath, func(p string, d os.FileInfo, err error) error {
		// Check if there was an error while walking.
		if err != nil {
			return err
		}

		// Skip files.
		if !d.IsDir() {
			return nil
		}

		// Convert the file path.
		targetPath, err := pathConverter(p)
		if err != nil {
			return err
		}

		// Add the file to the list.
		*directories = append(*directories, targetPath)
		return nil
	})
}

func listFiles(fs afero.Fs, sourcePath string, files *[]string, pathConverter func(string) (string, error)) error {
	return afero.Walk(fs, sourcePath, func(p string, d os.FileInfo, err error) error {
		// Check if there was an error while walking.
		if err != nil {
			return err
		}

		// Skip directories.
		if d.IsDir() {
			return nil
		}

		// Convert the file path.
		targetPath, err := pathConverter(p)
		if err != nil {
			return err
		}

		// Add the file to the list.
		*files = append(*files, targetPath)
		return nil
	})
}

func exclude(s string, fns []func(string) bool) bool {
	for _, fn := range fns {
		if fn(s) {
			return true
		}
	}

	return false
}

func excludeExactly(s string) func(string) bool {
	return func(p string) bool {
		return p == s
	}
}

func excludePrefix(s string) func(string) bool {
	return func(p string) bool {
		return strings.HasPrefix(p, s)
	}
}

func appendPathSeparator(s string) string {
	return s + string(filepath.Separator)
}
