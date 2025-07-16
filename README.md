# MCP Bins Server

A Model Context Protocol (MCP) server written in Go using the [mcp-go](https://github.com/mark3labs/mcp-go) library.

## Features

- **Bin Collection Tool**: Get bin collection dates for Reading addresses using UPRN

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
go build
```

## Testing

Run the test suite to verify functionality:

```bash
go test -v
```

The tests cover:

- Bin collection tool with various argument combinations
- HTTP client mocking and error handling
- API response validation and edge cases

## Usage

### Standalone

The server communicates via stdio transport:

```bash
./mcp-bins
```

You can optionally set a default UPRN for bin collection queries:

```bash
./mcp-bins -uprn 000000000
```

### Claude Desktop Integration

To use this MCP server with Claude Desktop:

1. Build the server:

   ```bash
   go build
   ```

2. Add the server to your Claude Desktop configuration. On macOS, edit:

   ```
   ~/Library/Application Support/Claude/claude_desktop_config.json
   ```

3. Add the following configuration:

   ```json
   {
     "mcpServers": {
       "mcp-bins": {
         "command": "/path/to/your/project/mcp-bins",
         "args": []
       }
     }
   }
   ```

   Or with a default UPRN:

   ```json
   {
     "mcpServers": {
       "mcp-bins": {
         "command": "/path/to/your/project/mcp-bins",
         "args": ["-uprn", "310045409"]
       }
     }
   }
   ```

4. Replace `/path/to/your/project/` with the actual path to your built binary

5. Restart Claude Desktop

6. You can now use the bin-collection tool in your conversations with Claude

**Note**: The bin-collection tool requires a valid UPRN (Unique Property Reference Number) for Reading addresses. You can find your UPRN on council tax bills or by searching the Reading Borough Council website.

## Tools

### bin-collection

- **Description**: Get bin collection dates for a Reading address using UPRN
- **Arguments**:
  - `uprn` (optional string): Unique Property Reference Number for the address
  - If no UPRN is provided, uses the default UPRN set at server startup
- **Example**: `bin-collection uprn=310045409`

## Dependencies

- Go 1.24.5+
- [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) v0.34.0
