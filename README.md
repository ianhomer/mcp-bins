# MCP Bins Server

A Model Context Protocol (MCP) server written in Go using the [mcp-go](https://github.com/mark3labs/mcp-go) library.

## Features

- **Bin Collection Tool**: Get bin collection dates for Reading Borough Council addresses using UPRN
- **Read-only**: Currently only supports reading collection schedules (no modifications)

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

You MUST set the default UPRN for bin collection queries. You can get it from
`curl https://api.reading.gov.uk/rbc/getaddresses/{postcode}`

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
         "args": ["-uprn", "000000000"]
       }
     }
   }
   ```

4. Replace `/path/to/your/project/` with the actual path to your built binary

5. Restart Claude Desktop

6. You can now use the bin-collection tool in your conversations with Claude

**Note**: The bin-collection tool requires a valid UPRN (Unique Property Reference Number) for Reading Borough Council addresses. You can find your UPRN on council tax bills or by searching the Reading Borough Council website. This tool provides read-only access to collection schedules.

## Tools

### bin-collection

- **Description**: Get bin collection dates for a Reading Borough Council address using UPRN (read-only)
- **Arguments**:
  - `uprn` (optional string): Unique Property Reference Number for the address
  - If no UPRN is provided, uses the default UPRN set at server startup
- **Example**: `bin-collection uprn=000000000`
- **Note**: This tool only reads collection schedules and cannot modify or update bin collection dates
- **Color Coding**:
  - ⚫ Black bin: Household/Domestic waste (general rubbish)
  - 🔴 Red bin: Recycling
  - 🟢 Green bin: Garden waste
- **Time Alerts**:
  - If it's 7AM-9AM on collection day: "⚠️ Collection is soon (around 9AM)!"
  - If it's after 9AM on collection day: "⚠️ Collection may have already happened (around 9AM)!"

## Dependencies

- Go 1.24.5+
- [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) v0.34.0
