package main

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

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
			setDefaultUPRN: "000000000",
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
			expectedOutput: "Upcoming bin collections for UPRN 000000000:\n\n📅 05/02/2020 00:00:00 (Wednesday)\n   🔴 Recycling Collection Service (red bin)\n\n",
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
						"uprn": "000000000",
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
						"uprn": "000000000",
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
			expectedOutput: "Upcoming bin collections for UPRN 310045409:\n\n📅 05/02/2020 00:00:00 (Wednesday)\n   🔴 Recycling Collection Service (red bin)\n\n📅 12/02/2020 00:00:00 (Wednesday)\n   ⚫ Household Waste Collection Service (black bin)\n\n",
		},
		{
			name: "successful response with no collections",
			request: mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name: "bin-collection",
					Arguments: map[string]interface{}{
						"uprn": "000000000",
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
			expectedOutput: "No upcoming bin collections found for UPRN 000000000",
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

func TestGetTimeAlert(t *testing.T) {
	// Test with a date that's today
	today := time.Now()
	todayStr := today.Format("02/01/2006 15:04:05")

	tests := []struct {
		name     string
		date     string
		hour     int
		expected string
	}{
		{
			name:     "early morning - no alert",
			date:     todayStr,
			hour:     6,
			expected: "",
		},
		{
			name:     "7AM - collection soon",
			date:     todayStr,
			hour:     7,
			expected: " ⚠️ Collection is soon (around 9AM)!",
		},
		{
			name:     "8AM - collection soon",
			date:     todayStr,
			hour:     8,
			expected: " ⚠️ Collection is soon (around 9AM)!",
		},
		{
			name:     "9AM - may have missed",
			date:     todayStr,
			hour:     9,
			expected: " ⚠️ Collection may have already happened (around 9AM)!",
		},
		{
			name:     "10AM - may have missed",
			date:     todayStr,
			hour:     10,
			expected: " ⚠️ Collection may have already happened (around 9AM)!",
		},
		{
			name:     "future date - no alert",
			date:     "05/02/2025 00:00:00",
			hour:     8,
			expected: "",
		},
		{
			name:     "past date - no alert",
			date:     "05/02/2020 00:00:00",
			hour:     8,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTime := time.Date(today.Year(), today.Month(), today.Day(), tt.hour, 0, 0, 0, time.Local)

			result := getTimeAlertWithTime(tt.date, mockTime)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
