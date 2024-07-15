package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jeremybower/tmpl"
	"github.com/spf13/cobra"
)

// Set with flags in the Makefile
var Version string

var configPaths *[]string
var outPath *string

var rootCmd = &cobra.Command{
	Use:     "tmpl [flags] template1 template2...",
	Example: "tmpl --config c1.yml --config c2.yml --out dest t1.tmpl t2.tmpl...",
	Args:    cobra.MinimumNArgs(1),
	Version: strings.TrimPrefix(Version, "v"),
	Short:   "Generates text from Go-style templates",
	Run: func(cmd *cobra.Command, args []string) {
		// Start a timer.
		start := time.Now()

		// Parse the templates
		t, err := tmpl.NewTemplate(args)
		exitIfError(cmd, err)

		// Check for an output path.
		if outPath == nil || *outPath == "" {
			exitWithUsageMessage(cmd, "missing required flag: --out")
		}

		// Create the output file.
		out, err := os.Create(*outPath)
		exitIfError(cmd, err)

		// Load the data.
		var data *tmpl.Data
		for _, path := range *configPaths {
			if data == nil {
				data, err = tmpl.NewData(path)
				exitIfError(cmd, err)
			} else {
				d, err := tmpl.NewData(path)
				exitIfError(cmd, err)
				data.Merge(d)
			}
		}

		// Execute the template.
		err = t.Execute(out, data.Config)
		exitIfError(cmd, err)

		// Close the file.
		err = out.Close()
		exitIfError(cmd, err)

		// Log stats
		fmt.Printf("Generated 1 file in %s\n", time.Since(start))
	},
}

func init() {
	outPath = rootCmd.Flags().StringP("out", "o", "", "path where the generated files will be written")
	configPaths = rootCmd.Flags().StringArrayP("config", "c", []string{}, "paths to configuration files")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func exitIfError(cmd *cobra.Command, err error) {
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}

func exitWithUsageMessage(cmd *cobra.Command, msg string) {
	cmd.PrintErrln(msg)
	cmd.Usage()
	os.Exit(2)
}
