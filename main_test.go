package main

import (
	"context"
	"io"
	"net/http"
	"strings"
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

// MockHTTPClient implements HTTPClient interface for testing
type MockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.response, m.err
}

func TestHandleBinCollection(t *testing.T) {
	ctx := context.Background()

	// Save original defaultUPRN and restore after test
	originalDefaultUPRN := defaultUPRN
	defer func() { defaultUPRN = originalDefaultUPRN }()

	tests := []struct {
		name           string
		request        mcp.CallToolRequest
		mockResponse   *http.Response
		mockError      error
		expectError    bool
		errorMsg       string
		expectedOutput string
		setDefaultUPRN string
	}{
		{
			name: "missing uprn without default",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "bin-collection",
					Arguments: map[string]interface{}{},
				},
			},
			expectError: true,
			errorMsg:    "uprn argument is required",
		},
		{
			name: "missing uprn with default",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "bin-collection",
					Arguments: map[string]interface{}{},
				},
			},
			setDefaultUPRN: "310045409",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"Collections": [
						{
							"Date": "05/02/2020 00:00:00",
							"Day": "Wednesday",
							"Service": "Recycling Collection Service"
						}
					]
				}`)),
			},
			expectError:    false,
			expectedOutput: "Upcoming bin collections for UPRN 310045409:\n\n📅 05/02/2020 00:00:00 (Wednesday)\n   🗑️ Recycling Collection Service\n\n",
		},
		{
			name: "invalid uprn format",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bin-collection",
					Arguments: map[string]interface{}{
						"uprn": "invalid",
					},
				},
			},
			expectError: true,
			errorMsg:    "uprn must be a valid number",
		},
		{
			name: "http client error",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bin-collection",
					Arguments: map[string]interface{}{
						"uprn": "310045409",
					},
				},
			},
			mockError:   http.ErrHandlerTimeout,
			expectError: true,
			errorMsg:    "failed to fetch bin collection data",
		},
		{
			name: "API returns error status",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bin-collection",
					Arguments: map[string]interface{}{
						"uprn": "310045409",
					},
				},
			},
			mockResponse: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("")),
			},
			expectError: true,
			errorMsg:    "API request failed with status 404",
		},
		{
			name: "successful response with collections",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bin-collection",
					Arguments: map[string]interface{}{
						"uprn": "310045409",
					},
				},
			},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"Collections": [
						{
							"Date": "05/02/2020 00:00:00",
							"Day": "Wednesday",
							"Service": "Recycling Collection Service"
						},
						{
							"Date": "12/02/2020 00:00:00",
							"Day": "Wednesday",
							"Service": "Household Waste Collection Service"
						}
					]
				}`)),
			},
			expectError:    false,
			expectedOutput: "Upcoming bin collections for UPRN 310045409:\n\n📅 05/02/2020 00:00:00 (Wednesday)\n   🗑️ Recycling Collection Service\n\n📅 12/02/2020 00:00:00 (Wednesday)\n   🗑️ Household Waste Collection Service\n\n",
		},
		{
			name: "successful response with no collections",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bin-collection",
					Arguments: map[string]interface{}{
						"uprn": "310045409",
					},
				},
			},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"Collections": []
				}`)),
			},
			expectError:    false,
			expectedOutput: "No upcoming bin collections found for UPRN 310045409",
		},
		{
			name: "invalid JSON response",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bin-collection",
					Arguments: map[string]interface{}{
						"uprn": "310045409",
					},
				},
			},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`invalid json`)),
			},
			expectError: true,
			errorMsg:    "failed to decode API response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set default UPRN for this test
			if tt.setDefaultUPRN != "" {
				defaultUPRN = tt.setDefaultUPRN
			} else {
				defaultUPRN = ""
			}

			mockClient := &MockHTTPClient{
				response: tt.mockResponse,
				err:      tt.mockError,
			}

			result, err := handleBinCollectionWithClient(ctx, tt.request, mockClient)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" {
					if !strings.Contains(err.Error(), tt.errorMsg) {
						t.Errorf("expected error to contain %q, got %q", tt.errorMsg, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else {
					if len(result.Content) != 1 {
						t.Errorf("expected 1 content item, got %d", len(result.Content))
					}
					textContent, ok := result.Content[0].(mcp.TextContent)
					if !ok {
						t.Errorf("expected TextContent, got %T", result.Content[0])
					}
					if textContent.Text != tt.expectedOutput {
						t.Errorf("expected output %q, got %q", tt.expectedOutput, textContent.Text)
					}
				}
			}
		})
	}
}
