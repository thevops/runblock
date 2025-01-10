# runblock

Run code blocks from Markdown files

## Description

`runblock` is a command-line tool that allows you to extract and execute code blocks from Markdown files.

## Features

- **List Code Blocks**: List all named code blocks from a Markdown file.
- **Print Code Blocks**: Print the content of a specific code block by name.
- **Execute Code Blocks**: Execute a specific code block by name in the current shell.

## Installation

To install `runblock`, you need to have Go installed on your machine. Then, you can clone the repository and build the project:

```sh
git clone https://github.com/yourusername/runblock.git
cd runblock
go build -o runblock
```

## Usage

### List Code Blocks

To list all named code blocks from a Markdown file, use the `list` command:

```
./runblock list --file path/to/yourfile.md
```

### Print Code Blocks

To print the content of a specific code block by name, use the `print` command:

```
./runblock print --file path/to/yourfile.md --name blockname
```

You can also include details about the code block:

```
./runblock print --file path/to/yourfile.md --name blockname --details
```

### Execute Code Blocks

To execute a specific code block by name in the current shell, use the `exec` command:

```
./runblock exec --file path/to/yourfile.md --name blockname
```

## Example

Given the following Markdown file [example/example.md](example/example.md):

You can list the code blocks:

```
./runblock list --file example/example.md
```

Output:

```
test1
test2 - This is a test 2
```

You can print a specific code block:

```
./runblock print --file example/example.md --name test1
```

Output:

```
echo "Hello test1"
```

You can execute a specific code block:

```
./runblock exec --file example/example.md --name test1
```

Output:

```
Hello test1
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## Acknowledgements

- [Cobra](https://github.com/spf13/cobra) for the CLI framework.
- [Goldmark](https://github.com/yuin/goldmark) for the Markdown parser.
- [Charmbracelet Log](https://github.com/charmbracelet/log) for the logging library.
