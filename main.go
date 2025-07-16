package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var defaultUPRN string

func main() {
	flag.StringVar(&defaultUPRN, "uprn", "", "Default UPRN for bin collection queries")
	flag.Parse()

	s := server.NewMCPServer(
		"mcp-bins",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Add bin collection tool
	var uprnDescription string
	if defaultUPRN != "" {
		uprnDescription = fmt.Sprintf("Unique Property Reference Number (UPRN) for the address (default: %s)", defaultUPRN)
	} else {
		uprnDescription = "Unique Property Reference Number (UPRN) for the address"
	}

	s.AddTool(
		mcp.NewTool("bin-collection",
			mcp.WithDescription("Get bin collection dates for a Reading address using UPRN"),
			mcp.WithString("uprn",
				mcp.Description(uprnDescription),
			),
		),
		handleBinCollection,
	)

	if err := server.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}

// BinCollection represents a single bin collection
type BinCollection struct {
	Date    string `json:"Date"`
	Day     string `json:"Day"`
	Service string `json:"Service"`
}

// BinCollectionResponse represents the API response
type BinCollectionResponse struct {
	Collections []BinCollection `json:"Collections"`
}

// HTTPClient interface for testing
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// Default HTTP client
var defaultHTTPClient HTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

func handleBinCollection(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return handleBinCollectionWithClient(ctx, request, defaultHTTPClient)
}

func handleBinCollectionWithClient(ctx context.Context, request mcp.CallToolRequest, client HTTPClient) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	uprnStr, ok := arguments["uprn"].(string)
	if !ok || uprnStr == "" {
		if defaultUPRN != "" {
			uprnStr = defaultUPRN
		} else {
			return nil, fmt.Errorf("uprn argument is required")
		}
	}

	// Validate UPRN is a number
	uprn, err := strconv.ParseInt(uprnStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("uprn must be a valid number: %v", err)
	}

	// Make API request
	url := fmt.Sprintf("https://api.reading.gov.uk/rbc/mycollections/%d", uprn)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bin collection data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var binData BinCollectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&binData); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %v", err)
	}

	// Format the response
	if len(binData.Collections) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("No upcoming bin collections found for UPRN %s", uprnStr),
				},
			},
		}, nil
	}

	var result string
	result = fmt.Sprintf("Upcoming bin collections for UPRN %s:\n\n", uprnStr)

	for _, collection := range binData.Collections {
		result += fmt.Sprintf("📅 %s (%s)\n   🗑️ %s\n\n", collection.Date, collection.Day, collection.Service)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}
