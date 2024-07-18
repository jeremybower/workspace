package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jeremybower/tmpl/internal"
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
					Usage:   "path(s) to configuration files",
				},
				&cli.StringFlag{
					Name:    "out",
					Aliases: []string{"o"},
					Usage:   "path to write the generated file",
				},
			},
			Action: func(c *cli.Context) error {
				// Start a timer.
				start := time.Now()

				// Check for the out flag.
				if !c.IsSet("out") {
					fmt.Fprintln(os.Stderr, "Error: The --out flag is required.")
					os.Exit(1)
				}

				// Parse the templates
				t, err := internal.NewTemplate(c.Args().Slice())
				exitIfError(err)

				// Create the output file.
				out, err := os.Create(c.String("out"))
				exitIfError(err)

				// Load the data.
				var data *internal.Data
				for _, path := range c.StringSlice("config") {
					if data == nil {
						data, err = internal.NewData(path)
						exitIfError(err)
					} else {
						d, err := internal.NewData(path)
						exitIfError(err)
						data.Merge(d)
					}
				}

				// Execute the template.
				err = t.Execute(out, data.Config)
				exitIfError(err)

				// Close the file.
				err = out.Close()
				exitIfError(err)

				// Log stats
				fmt.Printf("Generated 1 file in %s\n", time.Since(start))
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
