# Tmlp

[![Release](https://img.shields.io/github/release/jeremybower/tmpl.svg)](https://github.com/jeremybower/tmpl/releases)
[![Go Report](https://goreportcard.com/badge/github.com/jeremybower/tmpl)](https://goreportcard.com/report/github.com/jeremybower/tmpl)

Tmpl is a command line tool that generates text from [Go templates](https://pkg.go.dev/text/template) and YAML configuration files.

It is a standalone tool that can be used with Go, Node.js, Python, Ruby, PHP, Rust, C++, or any other language or framework you are using. This is especially helpful if you are writing multiple services in different languages and want a consistent approach when generating text files.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
  - [Generating a Dockerfile](#generating-a-dockerfile)
- [Contributing](#contributing)

## Features

- Supports any text-based format (Dockerfile, etc.)
- Available for Linux, macOS and Windows
- Lightweight and efficient
- Customizable with one or more config files
- Include multiple template files with glob filepaths

## Installation

**macOS**

Install the binary directly on macOS for `amd64`:

```sh
$ sudo curl -fsSL -o /usr/local/bin/tmpl https://github.com/jeremybower/tmpl/releases/latest/download/tmpl-darwin-amd64
$ sudo chmod +x /usr/local/bin/tmpl
```

or `arm64`:

```sh
$ sudo curl -fsSL -o /usr/local/bin/tmpl https://github.com/jeremybower/tmpl/releases/latest/download/tmpl-darwin-arm64
$ sudo chmod +x /usr/local/bin/tmpl
```

**Linux**

Install the binary directly on Linux for `amd64`:

```sh
$ sudo curl -fsSL -o /usr/local/bin/tmpl https://github.com/jeremybower/tmpl/releases/latest/download/tmpl-linux-amd64
$ sudo chmod +x /usr/local/bin/tmpl
```

or `arm64`:

```sh
$ sudo curl -fsSL -o /usr/local/bin/tmpl https://github.com/jeremybower/tmpl/releases/latest/download/tmpl-linux-arm64
$ sudo chmod +x /usr/local/bin/tmpl
```

**Windows**

Install the binary directly on Windows for `amd64`:

```ps
PS> Invoke-WebRequest -Uri 'https://github.com/jeremybower/tmpl/releases/latest/download/tmpl-windows-amd64.exe' -OutFile 'c:\temp\tmpl.exe'
```

## Usage

```sh
$ tmpl --help
Generates text from Go-style templates

Usage:
  tmpl [flags] template1 template2...

Examples:
tmpl --config c1.yml --config c2.yml --out dest t1.tmpl t2.tmpl...

Flags:
  -c, --config stringArray   paths to configuration files
  -h, --help                 help for tmpl
  -o, --out string           path where the generated files will be written
  -v, --version              version for tmpl
```

Tmpl accepts multiple config files, a single destination file and multiple templates. For example, generating a Dockerfile might require top-level configuration and local configuration files, plus other templates:

```sh
tpml -c ../config.yml -c Dockerfile.yml -o Dockerfile Dockerfile.tmpl includes/*.tmpl
```

See [Generating a Dockerfile](#generating-a-dockerfile) for the complete example.

## Template Functions

Tmpl includes all the functions provided by [sprig](http://masterminds.github.io/sprig/) and additional functions that support working with multiple templates and config files:

| Function        | Description                                                                                                                                                                                                                                                                                       |
| --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `globFilter`    | Accepts a glob pattern as the first parameter and a list of files as the second parameter. It is useful when combined with `listTemplates` to filter the set of templates to just those matching a specific pattern. For example, to include all templates at a specific path (`include/*.tmpl`). |
| `include`       | Similar to the standard `template` function, but the first parameter is a pipeline. The second parameter is the data to pass to the named template.                                                                                                                                               |
| `listTemplates` | Accepts no parameters and returns the sorted list of paths to all templates provided as arguments to the `tmpl` command.                                                                                                                                                                          |
| `require`       | Similar to the standard substitution approach of `{{ .Name }}`, but requires that the value is not the zero value. For example, `{{ require .Name }}`.                                                                                                                                            |

## Examples

### Generating a Dockerfile

The full example is available in `examples/dockerfile`:

```txt
.
├── config.yml
└── greeting
    ├── Dockerfile <--- generate this file
    ├── Dockerfile.tmpl
    ├── Dockerfile.yml
    ├── Makefile
    └── includes
        ├── en.tmpl
        └── fr.tmpl
```

The Dockerfile can be generated with:

```sh
cd examples/dockerfile/greeting
make Dockerfile
make build
make clean
```

There are three parts to this example: config files, template files and the command to generate.

#### Step 1: Config Files

Setup a top-level configuration file to use the same defaults for all Dockerfiles in a project:

`config.yml`:

```yml
Config:
  LanguageCode: "en"
  BaseImage: "ubuntu:24.04"
```

Use a local configuration file to change the generated output:

`greeting/Dockerfile.yml`:

```yml
Config:
  LanguageCode: "fr"
```

Tmpl will combine the configuration files in order where successive configuration files can overwrite earlier ones. Internally, templ will use this configuration:

```yml
Config:
  LanguageCode: "fr"
  BaseImage: "ubuntu:24.04"
```

Notice that `LanguageCode` was overwritten by the second configuration file.

#### Step 2: Templates

Create a main template for the Dockerfile:

`greeting/Dockerfile.tmpl`:

```Dockerfile
FROM {{ .BaseImage }}

{{ include (printf "includes/%s.tmpl" .LanguageCode) }}
```

Create template files to include for each language:

`greeting/includes/en.tmpl`:

```Dockerfile
CMD ["echo", "Hello!"]
```

`greeting/includes/fr.tmpl`:

```Dockerfile
CMD ["echo", "Bonjour!"]
```

#### Step 3: Generate

To generate the Dockerfile, run:

```sh
$ cd greeting
$ tpml -c ../config.yml,Dockerfile.yml -o Dockerfile Dockerfile.tmpl includes/*.tmpl
Generated 1 file in 2.999958ms
```

The resulting `Dockerfile` contains:

```Dockerfile
FROM ubuntu:24.04

CMD ["echo", "Bonjour!"]
```

## Contributing

Tmpl is written in Go. Pull requests are welcome.
