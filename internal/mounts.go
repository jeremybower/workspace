package internal

import (
	"errors"
	"os"
	"slices"

	"github.com/spf13/afero"
)

type Mounts []*Mount

func NewMounts(fs afero.Fs, specs []string) (Mounts, error) {
	// Load all the mounts.
	var mounts Mounts
	for _, spec := range specs {
		mount, err := NewMount(fs, spec)
		if err != nil {
			return nil, err
		}

		mounts = append(mounts, mount)
	}

	// Since the ealier mounts should take precedence over later mounts,
	// reverse the mounts.
	slices.Reverse(mounts)
	return mounts, nil
}

func (m Mounts) Directories(pattern string) ([]string, error) {
	// Iterate over the mounts and list the directories. The mounts ealier in the
	// list take precedence over the mounts later in the list.
	var directories []string
	var excludeFns []func(string) bool
	for _, mount := range m {
		err := mount.Dirs(pattern, &directories, excludeFns)
		if err != nil {
			return nil, err
		}

		// Since directories are pre-processed when mounted, a simple path matching
		// approach can be used to exclude directories with the same root.
		if mount.directory {
			excludeFns = append(excludeFns, excludePrefix(appendPathSeparator(mount.targetPath)))
		} else {
			excludeFns = append(excludeFns, excludeExactly(mount.targetPath))
		}
	}

	// Sort the directories if more than one mount to maintain predictable
	// results when multiple mounts match the same pattern.
	if len(m) > 1 {
		slices.Sort(directories)
	}

	return directories, nil
}

func (m Mounts) Files(pattern string) ([]string, error) {
	// Iterate over the mounts and list the files. The mounts ealier in the list
	// take precedence over the mounts later in the list.
	var files []string
	var excludeFns []func(string) bool
	for _, mount := range m {
		err := mount.Files(pattern, &files, excludeFns)
		if err != nil {
			return nil, err
		}

		// Since files are pre-processed when mounted, a simple path matching
		// approach can be used to exclude directories with the same root.
		if mount.directory {
			excludeFns = append(excludeFns, excludePrefix(appendPathSeparator(mount.targetPath)))
		} else {
			excludeFns = append(excludeFns, excludeExactly(mount.targetPath))
		}
	}

	// Sort the files if more than one mount to maintain predictable results
	// when multiple mounts match the same pattern.
	if len(m) > 1 {
		slices.Sort(files)
	}

	return files, nil
}

func (m Mounts) ReadFileString(targetPath string) (string, error) {
	// Iterate over the mounts and read the file. The mounts ealier in the list
	// take precedence over the mounts later in the list.
	for _, mount := range m {
		s, err := mount.ReadFileString(targetPath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			return "", err
		}

		return s, nil
	}

	// File not found.
	return "", os.ErrNotExist
}
