package main

import (
	"fmt"
	"os"

	"github.com/PlusLemon/mcp-cobra/mcp"
	"github.com/spf13/cobra"
)

func main() {
	// originalStdout := os.Stdout
	// r, w, _ := os.Pipe()
	// os.Stdout = w

	rootCmd := &cobra.Command{
		Use:   "foo",
		Short: "Foo Demo CLI",
	}

	// 定义子命令
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

	// fullArgs := []string{
	// 	"greet",
	// 	"--name", "fubang",
	// }
	// rootCmd.SetArgs(fullArgs)
	// err := rootCmd.Execute()
	// w.Close()
	// var buf bytes.Buffer
	// io.Copy(&buf, r)
	// os.Stdout = originalStdout
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// capturedText := buf.String()
	// fmt.Println("捕获的输出: " + capturedText)
}
