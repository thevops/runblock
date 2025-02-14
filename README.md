# runblock

Run code blocks from Markdown files

## Description

`runblock` is a command-line tool that allows you to extract and execute code blocks from Markdown files.

## Features

- **List Code Blocks**: List all named code blocks from a Markdown file.
- **Print Code Blocks**: Print the content of a specific code block by name.
- **Execute Code Blocks**: Execute a specific code block by name in the current shell.
- **Interactive Mode**: Run `runblock` in interactive mode to select and execute code blocks from a list.

## Installation

### From release

Download the latest release:

```sh
wget -O runblock "https://github.com/thevops/runblock/releases/latest/download/runblock_$(uname -s | tr '[:upper:]' '[:lower:]')_$(uname -m)"
chmod +x runblock
```

### From source

To install `runblock`, you need to have Go installed on your machine. Then, you can clone the repository and build the project:

```sh
git clone https://github.com/yourusername/runblock.git
cd runblock
go build -o runblock
```

## Usage

Run `runblock --help` to see the available commands and options.

## Examples

Check out the [example](./example) directory for a sample Markdown file and code blocks.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## Acknowledgements

- [Cobra](https://github.com/spf13/cobra) for the CLI framework.
- [Goldmark](https://github.com/yuin/goldmark) for the Markdown parser.
- [Charmbracelet Log](https://github.com/charmbracelet/log) for the logging library.
