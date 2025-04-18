# oeN (Neo backwards)

A Matrix-inspired anti-agent CLI tool that enables interaction with Claude 3 through tools. Based on the architecture described in [How to Build an Agent](https://ampcode.com/how-to-build-an-agent) with additional features.

## Overview

oeN is a command-line interface for interacting with Claude 3.7 Sonnet, that implements the agent pattern enabling Claude to perform various file system operations through defined tools:

- `read_file`: Read the contents of a file
- `list_files`: List files in a directory
- `edit_file`: Make changes to a text file
- `make_directory`: Create new directories
- `remove_directory`: Remove directories (with optional recursive deletion)
- `rename_directory`: Rename or move directories

The tool acts as a bridge between Claude's reasoning capabilities and your local file system, allowing you to have Claude help with file management tasks through natural language.

## Installation

### Prerequisites

- Go
- Anthropic API key (set as environment variable)

### Building from source

```bash
git clone https://github.com/MarkusZoppelt/oen.git
cd oen
go build
```

Or install directly using Go:

```bash
go install github.com/MarkusZoppelt/oen@latest
```

## Usage

1. Set your Anthropic API key as an environment variable:
   ```bash
   export ANTHROPIC_API_KEY="your-api-key"
   ```

2. Run the application:
   ```bash
   ./oen
   ```

3. Start conversing with Claude through the CLI interface. You can ask it to perform file operations, and it will use the appropriate tools to execute them.

## Example Interactions

```
You: Can you list the files in the current directory?
Claude: I'll list the files in the current directory for you.

tool: list_files({})
[".git/", "cmd/", "go.mod", "go.sum", "main.go", "oen"]

Here are the files and directories in the current location:
- .git/ (directory)
- cmd/ (directory)
- go.mod
- go.sum
- main.go
- oen (executable)

You: Create a new directory called "test"
Claude: I'll create a new directory called "test" for you.

tool: make_directory({"path":"test"})
Successfully created directory test

I've created a new directory named "test" in the current location.

You: Create a file in that directory called hello.txt with the content "Hello, World!"
Claude: I'll create that file for you.

tool: edit_file({"path":"test/hello.txt","old_str":"","new_str":"Hello, World!"})
Successfully created file test/hello.txt

I've created a file named "hello.txt" in the "test" directory with the content "Hello, World!".
```

## Architecture

The oeN implementation follows the agent pattern where:

1. User input is collected through a CLI interface
2. Input is forwarded to Claude 3.7 Sonnet with tool definitions
3. Claude determines which tool to use based on the request
4. The tool execution is handled locally and results are returned to Claude
5. Claude provides a natural language response integrating the tool results

This creates a conversational experience where Claude can reason about and perform file system operations through defined tools.

## License

[MIT License](LICENSE)
