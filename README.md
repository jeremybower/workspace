# Tmlp

[![Release](https://img.shields.io/github/release/jeremybower/tmpl.svg)](https://github.com/jeremybower/tmpl/releases)
[![Go Report](https://goreportcard.com/badge/github.com/jeremybower/tmpl)](https://goreportcard.com/report/github.com/jeremybower/tmpl)

Tmpl is a command line tool that generates text from [Go templates](https://pkg.go.dev/text/template) and YAML configuration files.

It is a standalone tool that can be used with Go, Node.js, Python, Ruby, PHP, Rust, C++, or any other language or framework you are using. This is especially helpful if you are writing multiple services in different languages and want a consistent approach when generating text files.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Template Functions](#template-functions)
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
$ tmpl generate --help
NAME:
   tmpl generate - Generate text from template and configuration files

USAGE:
   tmpl generate [command options]

OPTIONS:
   --config value, -c value [ --config value, -c value ]  Apply configuration data to the templates
   --mount value, -m value [ --mount value, -m value ]    Attach a filesystem mount to the template engine
   --out value, -o value                                  Write the generated text to file
   --help, -h                                             show help
```

Tmpl accepts multiple config files, a single destination file and mounts to access the file system at known paths. For example, generating a Dockerfile might require general configuration and specific configuration files, plus other templates:

```sh
tmpl generate -c config.yml -c Dockerfile.yml -m Dockerfile.tmpl:/Dockerfile -m includes:/includes -o Dockerfile /Dockerfile.tmpl
```

See [Generating a Dockerfile](#generating-a-dockerfile) for the complete example.

## Template Functions

Tmpl includes all the functions provided by [sprig](http://masterminds.github.io/sprig/) and additional functions that support working with multiple templates and config files:

| Function      | Description                                                                                                                                                                              |
| ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `dirs`        | Lists all the directories that were mounted. The only parameter is a glob pattern to match against the directory names.                                                                  |
| `filename`    | Returns the filename of the current template.                                                                                                                                            |
| `files`       | Lists all the files that were mounted. The only parameter is a glob pattern to match against the file names.                                                                             |
| `include`     | Similar to the standard `template` function, but the first parameter accepts a pipeline to select templates dynamically. The second parameter is the data to pass to the named template. |
| `includeText` | Similar to `include` function, but passes the file's text through unchanged. The only parameter is a pipeline to select the files dynamically.                                           |

## Examples

### Generating a Dockerfile

The full example is available in `examples/dockerfile`:

```txt
.
├── Dockerfile <--- generate this file
├── Dockerfile.tmpl
├── Dockerfile.yml
├── Makefile
├── config.yml
└── includes
    ├── en.tmpl
    └── fr.tmpl
```

The Dockerfile can be generated with:

```sh
cd examples/dockerfile
make Dockerfile
make build
make clean
```

There are three parts to this example: config files, template files and the command to generate.

#### Step 1: Config Files

Setup a general configuration file to use the same defaults for all Dockerfiles in a project:

`config.yml`:

```yml
Config:
  LanguageCode: "en"
  BaseImage: "ubuntu:24.04"
```

Use a specific configuration file to change the generated output:

`Dockerfile.yml`:

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

`Dockerfile.tmpl`:

```Dockerfile
FROM {{ .BaseImage }}

{{ include (printf "includes/%s.tmpl" .LanguageCode) . }}
```

Create template files to include for each language:

`includes/en.tmpl`:

```Dockerfile
CMD ["echo", "Hello!"]
```

`includes/fr.tmpl`:

```Dockerfile
CMD ["echo", "Bonjour!"]
```

#### Step 3: Generate

To generate the Dockerfile, run:

```sh
$ make Dockerfile
Generated 1 file(s) in 6.961916ms
/tmpl/examples/dockerfile/Dockerfile
```

The resulting `Dockerfile` contains:

```Dockerfile
FROM ubuntu:24.04

CMD ["echo", "Bonjour!"]
```

## Contributing

Tmpl is written in Go. Pull requests are welcome.
