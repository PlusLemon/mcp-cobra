package mcp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type MCPServer struct {
	rootCmd *cobra.Command
	server  *server.MCPServer
}

func NewMCPServer(rootCmd *cobra.Command) *MCPServer {
	s := &MCPServer{
		rootCmd: rootCmd,
		server: server.NewMCPServer(rootCmd.Short, "1.0.0",
			server.WithLogging(),
			server.WithRecovery(),
			server.WithResourceCapabilities(true, true),
		),
	}

	leaves := getLeafCommands(rootCmd)

	for _, leaf := range leaves {
		fullPath := getFullCommandPath(leaf)
		toolName := strings.Join(fullPath, " ")
		// toolName := leaf.Name()
		toolDesc := leaf.Short
		if toolDesc == "" {
			toolDesc = leaf.Long
		}
		var toolOptions []mcp.ToolOption
		toolOptions = append(toolOptions, mcp.WithDescription(toolDesc))

		flags := getAllFlagDefs(leaf)
		for _, f := range flags {
			var propertieOptions []mcp.PropertyOption
			if f.DefValue == "" {
				propertieOptions = append(propertieOptions, mcp.Required())
			}
			propertieOptions = append(propertieOptions, mcp.Description(f.Usage))
			switch f.Value.Type() {
			case "string":
				propertieOptions = append(propertieOptions, mcp.DefaultString(f.DefValue))
				toolOptions = append(toolOptions, mcp.WithString(f.Name, propertieOptions...))
			case "int":
				defaultValue, _ := strconv.ParseFloat(f.DefValue, 64)
				propertieOptions = append(propertieOptions, mcp.DefaultNumber(defaultValue))
				toolOptions = append(toolOptions, mcp.WithNumber(f.Name, propertieOptions...))
			case "bool":
				defaultValue, _ := strconv.ParseBool(f.DefValue)
				propertieOptions = append(propertieOptions, mcp.DefaultBool(defaultValue))
				toolOptions = append(toolOptions, mcp.WithBoolean(f.Name, propertieOptions...))
			case "float32", "float64":
				defaultValue, _ := strconv.ParseFloat(f.DefValue, 64)
				propertieOptions = append(propertieOptions, mcp.DefaultNumber(defaultValue))
				toolOptions = append(toolOptions, mcp.WithNumber(f.Name, propertieOptions...))
			default:
				toolOptions = append(toolOptions, mcp.WithString(f.Name, propertieOptions...))
			}
		}
		tool := mcp.NewTool(toolName, toolOptions...)

		LogInfo(fmt.Sprintf("Registering tool %s with description: %s", toolName, toolDesc))
		s.server.AddTool(tool, s.handleToolCall(leaf))
	}

	return s
}

func (s *MCPServer) ServeStdio() error {
	return server.ServeStdio(s.server)
}

func (s *MCPServer) handleToolCall(cmd *cobra.Command) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		LogInfo(fmt.Sprintf("Invoking tool %s with parameters: %+v", cmd.Name(), request.Params.Arguments))

		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		fullArgs := []string{}
		for _, part := range getFullCommandPath(cmd) {
			fullArgs = append(fullArgs, part)
		}

		for key, val := range request.Params.Arguments {
			if key == "args" {
				continue
			}
			fullArgs = append(fullArgs, "--"+key)
			fullArgs = append(fullArgs, fmt.Sprintf("%v", val))
		}

		if args, ok := request.Params.Arguments["args"].([]interface{}); ok {
			for _, arg := range args {
				fullArgs = append(fullArgs, fmt.Sprintf("%v", arg))
			}
		}

		LogInfo(fmt.Sprintf("Executing command: %v\n", fullArgs))
		s.rootCmd.SetArgs(fullArgs)

		err := s.rootCmd.Execute()

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		os.Stdout = originalStdout
		if err != nil {
			LogInfo(fmt.Sprintf("Error executing command: %v", err))
			return mcp.NewToolResultText(err.Error()), nil
		}
		capturedText := buf.String()
		return mcp.NewToolResultText(capturedText), nil
	}
}

func getLeafCommands(cmd *cobra.Command) []*cobra.Command {
	var leaves []*cobra.Command
	if len(cmd.Commands()) == 0 {
		leaves = append(leaves, cmd)
	} else {
		for _, sub := range cmd.Commands() {
			leaves = append(leaves, getLeafCommands(sub)...)
		}
	}
	return leaves
}

func getFullCommandPath(cmd *cobra.Command) []string {
	if cmd.Parent() == nil {
		// ignore the root command
		return []string{}
	}
	parentPath := getFullCommandPath(cmd.Parent())
	return append(parentPath, cmd.Name())
}

func getAllFlagDefs(cmd *cobra.Command) map[string]*pflag.Flag {
	flags := make(map[string]*pflag.Flag)
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		flags[f.Name] = f
	})
	cmd.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		flags[f.Name] = f
	})
	return flags
}
