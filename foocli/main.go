package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PlusLemon/mcp-cobra/mcp"
	"github.com/spf13/cobra"
)

func main() {
	defer mcp.CloseGlobalLogger()

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

	mcp.LogInfo(fmt.Sprintf("Starting foo CLI, args: %v", os.Args))
	if len(os.Args) > 1 && os.Args[len(os.Args)-1] == "mcp-server" {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		mcpServer := mcp.NewMCPServer(rootCmd)
		go func() {
			if err := mcpServer.ServeStdio(); err != nil {
				fmt.Printf("MCP server error: %v\n", err)
				os.Exit(1)
			}
		}()
		<-quit
	} else {
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
