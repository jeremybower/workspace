package internal

import (
	"os"
	"path"
	"time"

	"github.com/spf13/afero"
)

type Result struct {
	Filenames []string
	Duration  time.Duration
}

func Execute(tmplFilename string, mountSpecs []string, configFilenames []string, outFilename string, fs afero.Fs) (*Result, error) {
	// Start the timer.
	start := time.Now()

	// Create the mounts.
	mounts, err := NewMounts(fs, mountSpecs)
	if err != nil {
		return nil, err
	}

	// Create the config spec.
	configSpec, err := NewConfigSpec(fs, configFilenames)
	if err != nil {
		return nil, err
	}

	// Convert the out filename to an absolute path.
	if !path.IsAbs(outFilename) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		outFilename = path.Clean(path.Join(wd, outFilename))
	}

	// Execute the template.
	err = execute(tmplFilename, mounts, configSpec, outFilename, fs)
	if err != nil {
		return nil, err
	}

	// Return the result.
	return &Result{
		Filenames: []string{outFilename},
		Duration:  time.Since(start),
	}, nil
}

func execute(tmplFilename string, mounts Mounts, configSpec *ConfigSpec, outFilename string, fs afero.Fs) error {
	// Create the template cache.
	templateManager := NewTemplateCache(mounts)

	// Create the template.
	t, err := templateManager.Template(tmplFilename)
	if err != nil {
		return err
	}

	// Create the out file.
	outFile, err := fs.Create(outFilename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Execute the template.
	return t.Execute(outFile, mounts, configSpec.config)
}
