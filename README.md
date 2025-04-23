# MCP-Cobra

MCP-Cobra is a Go library that integrates the [Cobra](https://github.com/spf13/cobra) command-line interface framework with the [Model Context Protocol (MCP)](https://github.com/mark3labs/mcp-go). This library allows you to expose your Cobra CLI commands as MCP tools, making them accessible to AI assistants and other MCP clients.

## Features

- Automatically convert Cobra commands to MCP tools
- Support for various flag types (string, int, bool, float)
- Seamless integration with existing Cobra-based CLIs
- Run your CLI application in traditional command-line mode or as an MCP server

## Installation

```bash
go get github.com/PlusLemon/mcp-cobra
```

## Usage

To use MCP-Cobra, you need to:

1. Create a Cobra command structure as usual
2. Initialize an MCP server with your root command
3. Serve the MCP server via stdio when needed

Example:

```go
package main

import (
    "fmt"
    "os"

    "github.com/PlusLemon/mcp-cobra/mcp"
    "github.com/spf13/cobra"
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "myapp",
        Short: "My application description",
    }

    // Add your commands and flags...

    // Check if the application should run as an MCP server
    if len(os.Args) > 1 && os.Args[1] == "mcp-server" {
        mcpServer := mcp.NewMCPServer(rootCmd)
        if err := mcpServer.ServeStdio(); err != nil {
            fmt.Printf("MCP server error: %v\n", err)
            os.Exit(1)
        }
    } else {
        // Run as a normal CLI application
        if err := rootCmd.Execute(); err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
    }
}
```

## Example: foocli

The `foocli` directory contains a complete working example of MCP-Cobra in action.

### What it does

`foocli` is a simple CLI application with a `greet` command that prints a greeting message. It can run both as a traditional CLI and as an MCP server.

### Running the example

As a CLI:

```bash
# Build the example
go build -o bin/foo ./foocli

# Run as a regular CLI
./bin/foo greet --name "World"
# Output: Hello, World!
```

As an MCP server:

```bash
# Run as an MCP server
./bin/foo mcp-server
```

When running as an MCP server, the `greet` command is exposed as an MCP tool that can be invoked by MCP clients.

### Example Code

```go
package main

import (
	"fmt"
	"os"

	"github.com/PlusLemon/mcp-cobra/mcp"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "foo",
		Short: "Foo Demo CLI",
	}

	// Define subcommand
	var greetWord string
	greetCmd := &cobra.Command{
		Use:   "greet",
		Short: "Greet someone",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Hello, %s!\n", greetWord)
		},
	}
	greetCmd.Flags().StringVar(&greetWord, "name", "Foo", "Name to greet")

	rootCmd.AddCommand(greetCmd)

	if len(os.Args) > 1 && os.Args[1] == "mcp-server" {
		mcpServer := mcp.NewMCPServer(rootCmd)
		if err := mcpServer.ServeStdio(); err != nil {
			fmt.Printf("MCP server error: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
```

## How It Works

When running in MCP server mode, MCP-Cobra analyzes your Cobra command structure and converts each command to an MCP tool. Command flags are converted to tool parameters with appropriate types and metadata.

When an MCP client calls a tool, MCP-Cobra:

1. Maps the tool call to the appropriate Cobra command
2. Converts the tool parameters to command flags
3. Executes the command
4. Captures the command output
5. Returns the result to the MCP client

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.