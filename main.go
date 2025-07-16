package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"hello-world-server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	// Add a hello world tool
	s.AddTool(
		mcp.NewTool("hello",
			mcp.WithDescription("Says hello to the world or a specific name"),
			mcp.WithString("name",
				mcp.Description("Optional name to greet"),
			),
		),
		handleHello,
	)

	// Add a simple resource
	s.AddResource(
		mcp.NewResource("hello://world",
			"Hello World",
			mcp.WithMIMEType("text/plain"),
		),
		handleHelloResource,
	)

	if err := server.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}

func handleHello(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	name := "World"
	if n, ok := arguments["name"].(string); ok && n != "" {
		name = n
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Hello, %s! This content came from the MCP server", name),
			},
		},
	}, nil
}

func handleHelloResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/plain",
			Text:     "Hello, World! This is a simple MCP resource.",
		},
	}, nil
}
