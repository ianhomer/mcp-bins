# Hello World MCP Server

A simple Go-based Model Context Protocol (MCP) server that demonstrates basic MCP functionality.

## Project Structure

- `main.go` - Main server implementation with hello tool and resource
- `go.mod` - Go module dependencies
- `README.md` - Project documentation

## Build & Test Commands

```bash
# Build the server
go build -o hello-world-server

# Run tests (if any)
go test ./...

# Clean up dependencies
go mod tidy
```

## Development Notes

- Uses mcp-go library v0.34.0
- Implements stdio transport for MCP communication
- Includes one tool (`hello`) and one resource (`hello://world`)
- Server handlers follow mcp-go v0.34.0 API patterns
