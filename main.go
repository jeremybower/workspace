package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/jeremybower/tmpl/internal"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

// Set with flags in the Makefile
var Version string

//go:embed LICENSE
var license string

func main() {
	if err := newApp().Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "tmpl"
	app.Usage = "Generates text from Go templates and YAML configuration files."
	app.Version = strings.TrimPrefix(Version, "v")

	app.Commands = []*cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate text from template and configuration files",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:    "config",
					Aliases: []string{"c"},
					Usage:   "Apply configuration data to the templates",
				},
				&cli.StringFlag{
					Name:        "missingkey",
					Usage:       "Controls the behavior during execution if a map is indexed with a key that is not present in the map",
					DefaultText: "error",
					Value:       "error",
				},
				&cli.StringSliceFlag{
					Name:    "mount",
					Aliases: []string{"m"},
					Usage:   "Attach a filesystem mount to the template engine",
				},
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
					Usage:   "Write the generated text to file",
				},
			},
			Action: func(c *cli.Context) error {
				// Check for the out flag.
				if !c.IsSet("out") {
					exitWithMessage("Error: The --out flag is required.")
				}

				// Check fo exactly one argument.
				if c.NArg() != 1 {
					exitWithMessage("Error: Exactly one argument is required.")
				}

				// Collect the options.
				opts := internal.Options{
					MissingKey: c.String("missingkey"),
				}

				// Execute the template.
				fs := afero.NewOsFs()
				templateFilename := c.Args().First()
				mountSpecs := c.StringSlice("mount")
				configFilenames := c.StringSlice("config")
				outFilename := c.String("out")
				result, err := internal.Execute(fs, templateFilename, mountSpecs, configFilenames, outFilename, opts)
				exitIfError(err)

				// Print the results.
				fmt.Printf("Generated %d file(s) in %s\n", len(result.Filenames), result.Duration)
				for _, filename := range result.Filenames {
					fmt.Println(filename)
				}

				// Success.
				return nil
			},
		},
		{
			Name:  "license",
			Usage: "Prints the license",
			Action: func(c *cli.Context) error {
				// Always include a newline at the end.
				if license[len(license)-1] != '\n' {
					license += "\n"
				}

				// Print the license.
				fmt.Print(license)
				return nil
			},
		},
	}

	return app
}

func exitIfError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func exitWithMessage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
