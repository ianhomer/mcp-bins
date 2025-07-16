# Hello World MCP Server

A simple Model Context Protocol (MCP) server written in Go using the [mcp-go](https://github.com/mark3labs/mcp-go) library.

## Features

- **Hello Tool**: Says hello to the world or a specific name
- **Hello Resource**: A simple text resource accessible at `hello://world`

## Development Setup

1. Clone the repository
2. Install pre-commit framework:
   ```bash
   pip install pre-commit
   ```
3. Install the pre-commit hooks:
   ```bash
   pre-commit install
   ```

The pre-commit hooks will automatically run before each commit and include:
- `go fmt` - Code formatting
- `go vet` - Static analysis
- `go mod tidy` - Clean up module dependencies
- `go test` - All tests must pass
- General checks for trailing whitespace, large files, etc.

## Building

```bash
go build -o hello-world-server
```

## Testing

Run the test suite to verify functionality:

```bash
go test -v
```

The tests cover:
- Hello tool with various argument combinations
- Hello resource content validation
- Error handling and edge cases

## Usage

The server communicates via stdio transport:

```bash
./hello-world-server
```

## Tools

### hello
- **Description**: Says hello to the world or a specific name
- **Arguments**:
  - `name` (optional string): Name to greet

## Resources

### hello://world
- **Type**: text/plain
- **Description**: A simple hello world resource

## Dependencies

- Go 1.24.5+
- [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) v0.34.0
