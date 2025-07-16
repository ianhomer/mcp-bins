package main

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestHandleHello(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		request  mcp.CallToolRequest
		expected string
	}{
		{
			name: "hello with no name",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "hello",
					Arguments: map[string]interface{}{},
				},
			},
			expected: "Hello, World! This content came from the MCP server",
		},
		{
			name: "hello with name",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "hello",
					Arguments: map[string]interface{}{
						"name": "Alice",
					},
				},
			},
			expected: "Hello, Alice! This content came from the MCP server",
		},
		{
			name: "hello with empty name",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "hello",
					Arguments: map[string]interface{}{
						"name": "",
					},
				},
			},
			expected: "Hello, World! This content came from the MCP server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handleHello(ctx, tt.request)
			if err != nil {
				t.Fatalf("handleHello() error = %v", err)
			}

			if len(result.Content) != 1 {
				t.Fatalf("expected 1 content item, got %d", len(result.Content))
			}

			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Fatalf("expected TextContent, got %T", result.Content[0])
			}

			if textContent.Text != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, textContent.Text)
			}
		})
	}
}

func TestHandleHelloResource(t *testing.T) {
	ctx := context.Background()

	request := mcp.ReadResourceRequest{
		Params: mcp.ReadResourceParams{
			URI: "hello://world",
		},
	}

	result, err := handleHelloResource(ctx, request)
	if err != nil {
		t.Fatalf("handleHelloResource() error = %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 resource content, got %d", len(result))
	}

	textResource, ok := result[0].(mcp.TextResourceContents)
	if !ok {
		t.Fatalf("expected TextResourceContents, got %T", result[0])
	}

	expectedText := "Hello, World! This is a simple MCP resource."
	if textResource.Text != expectedText {
		t.Errorf("expected %q, got %q", expectedText, textResource.Text)
	}

	if textResource.URI != "hello://world" {
		t.Errorf("expected URI %q, got %q", "hello://world", textResource.URI)
	}

	if textResource.MIMEType != "text/plain" {
		t.Errorf("expected MIME type %q, got %q", "text/plain", textResource.MIMEType)
	}
}
