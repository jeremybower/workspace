package internal

import (
	"os"
	"path"
	"time"

	"github.com/spf13/afero"
)

type Options struct {
	MissingKey string
}

func DefaultOptions() Options {
	return Options{
		MissingKey: "error",
	}
}

type Result struct {
	Filenames []string
	Duration  time.Duration
}

func Execute(fs afero.Fs, tmplFilename string, mountSpecs []string, configFilenames []string, outFilename string, opts Options) (*Result, error) {
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
	err = execute(fs, tmplFilename, mounts, configSpec, outFilename, opts)
	if err != nil {
		return nil, err
	}

	// Return the result.
	return &Result{
		Filenames: []string{outFilename},
		Duration:  time.Since(start),
	}, nil
}

func execute(fs afero.Fs, tmplFilename string, mounts Mounts, configSpec *ConfigSpec, outFilename string, opts Options) error {
	// Create the template cache.
	templateManager := NewTemplateCache(mounts, opts)

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
