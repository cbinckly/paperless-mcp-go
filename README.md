# Paperless MCP Go

An MCP (Model Context Protocol) Server for [Paperless-ngx](https://docs.paperless-ngx.com/), enabling AI assistants and other MCP clients to interact with your Paperless document management system.

## Quickstart

Get up and running quickly...

### Run Your Own Instance

You're already running your own Paperless Server, this will be much easier...

Create an environment file:

```env
PAPERLESS_URL=https://paperless.example.com
PAPERLESS_TOKEN=your_paperless_api_token_here
MCP_AUTH_TOKEN=optional_mcp_auth_token
LOG_LEVEL=info
MCP_TRANSPORT=http
MCP_HTTP_PORT=8080
```

Run with docker:

```bash
$ docker run --rm --env-file .env -p 8080:8080  cbinckly/paperless-mcp-go:latest
```

And you're up and running on `http://localhost:8080/mcp`.

### Run with Claude or any MCP Client

Docker is still easiest, if you've got it, update your MCP config:

```json
{
  "mcpServers": {
    "paperless-mcp": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "--pull=always",
        "cbinckly/paperless-mcp-go:latest"
      ],
      "env": {
        "PAPERLESS_URL": "https://paperless.example.com",
        "PAPERLESS_TOKEN": "your_paperless_api_token_here",
        "MCP_TRANSPORT": "stdio"
      }
    }
  }
}
```

## Features

- **Complete Document Management**: Search, retrieve, create, update, and delete documents
- **Bulk Operations**: Efficiently edit multiple documents at once
- **Metadata Management**: Full CRUD operations for correspondents, document types, tags, storage paths, and custom fields
- **Dual Transport Modes**: Support for both stdio and HTTP Streaming transports
- **Docker Support**: Multi-stage Docker builds with security best practices
- **Comprehensive Logging**: Structured logging with configurable levels
- **Type-Safe**: Built with Go 1.23 for reliability and performance

### Available MCP Tools

#### Document Tools
- `search_documents` - Search for documents by text query with pagination
- `find_similar_documents` - Find documents similar to a given document
- `get_document` - Retrieve a document by ID with all metadata
- `get_document_content` - Get the text content of a document
- `create_document` - Create a new document
- `update_document` - Update document metadata
- `delete_document` - Delete a document
- `bulk_edit_documents` - Perform bulk operations on multiple documents

#### Correspondent Tools
- `list_correspondents` - List all correspondents with pagination
- `get_correspondent` - Get correspondent details by ID
- `create_correspondent` - Create a new correspondent
- `update_correspondent` - Update correspondent information
- `delete_correspondent` - Delete a correspondent

#### Document Type Tools
- `list_document_types` - List all document types with pagination
- `get_document_type` - Get document type details by ID
- `create_document_type` - Create a new document type
- `update_document_type` - Update document type information
- `delete_document_type` - Delete a document type

#### Tag Tools
- `list_tags` - List all tags with pagination
- `get_tag` - Get tag details by ID
- `create_tag` - Create a new tag
- `update_tag` - Update tag information
- `delete_tag` - Delete a tag

#### Storage Path Tools
- `list_storage_paths` - List all storage paths with pagination
- `get_storage_path` - Get storage path details by ID
- `create_storage_path` - Create a new storage path
- `update_storage_path` - Update storage path information
- `delete_storage_path` - Delete a storage path

#### Custom Field Tools
- `list_custom_fields` - List all custom fields with pagination
- `get_custom_field` - Get custom field details by ID
- `create_custom_field` - Create a new custom field
- `update_custom_field` - Update custom field information
- `delete_custom_field` - Delete a custom field

#### Utility Tools
- `ping` - Test tool that returns pong
- `server_info` - Get MCP server and Paperless connection information

## Prerequisites

- Go 1.23 or later (for building from source)
- A running [Paperless-ngx](https://docs.paperless-ngx.com/) instance
- Paperless-ngx API token
- Docker and Docker Compose (optional, for containerized deployment)

## Installation

### From Source

```bash
git clone https://github.com/cbinckly/paperless-mcp-go.git
cd paperless-mcp-go
go mod download
```

### Using Docker

```bash
git clone https://github.com/cbinckly/paperless-mcp-go.git
cd paperless-mcp-go
docker-compose up -d
```

## Configuration

The server is configured using environment variables. Create a `.env` file based on `.env.example`:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PAPERLESS_URL` | **Yes** | - | URL of your Paperless-ngx instance |
| `PAPERLESS_TOKEN` | **Yes** | - | API token for Paperless-ngx authentication |
| `MCP_AUTH_TOKEN` | No | - | Optional authentication token for MCP clients |
| `LOG_LEVEL` | No | `info` | Logging level: `debug`, `info`, `warn`, `error` |
| `MCP_TRANSPORT` | No | `stdio` | Transport mode: `stdio` or `http` |
| `MCP_HTTP_PORT` | No | `8080` | HTTP port (only used when `MCP_TRANSPORT=http`) |

### Example `.env` File

```env
PAPERLESS_URL=https://paperless.example.com
PAPERLESS_TOKEN=your_paperless_api_token_here
MCP_AUTH_TOKEN=optional_mcp_auth_token
LOG_LEVEL=info
MCP_TRANSPORT=http
MCP_HTTP_PORT=8080
```

## Building

### Build from Source

```bash
go build -o paperless-mcp ./cmd/server
```

### Build Docker Image

```bash
docker build -t paperless-mcp-go .
```

## Running

### Stdio Transport (Default)

Stdio transport is typically used when the MCP server is launched by an MCP client:

```bash
# Set environment variables
export PAPERLESS_URL=https://paperless.example.com
export PAPERLESS_TOKEN=your_token_here

# Run the server
./paperless-mcp
```

### HTTP Transport

HTTP transport runs the server as a standalone service using the modern StreamableHTTP protocol:

```bash
# Set environment variables
export PAPERLESS_URL=https://paperless.example.com
export PAPERLESS_TOKEN=your_token_here
export MCP_TRANSPORT=http
export MCP_HTTP_PORT=8080

# Run the server
./paperless-mcp
```

The server will be available at `http://localhost:8080`.

### Using Docker

```bash
docker run --env-file .env -p 8080:8080 paperless-mcp-go
```

### Using Docker Compose

```bash
docker-compose up -d
```

Check logs:
```bash
docker-compose logs -f
```

Stop the server:
```bash
docker-compose down
```

## Usage Examples

### Example 1: Search for Documents

Using the `search_documents` tool to find documents containing "invoice":

```json
{
  "tool": "search_documents",
  "arguments": {
    "query": "invoice",
    "page": 1,
    "page_size": 10
  }
}
```

### Example 2: Create a Tag

Using the `create_tag` tool to create a new tag:

```json
{
  "tool": "create_tag",
  "arguments": {
    "name": "Urgent",
    "color": "#ff0000"
  }
}
```

### Example 3: Bulk Edit Documents

Using the `bulk_edit_documents` tool to add tags to multiple documents:

```json
{
  "tool": "bulk_edit_documents",
  "arguments": {
    "document_ids": [1, 2, 3],
    "add_tags": [5, 6],
    "set_correspondent": 2
  }
}
```

## Deployment

### Production Considerations

1. **Security**:
   - Use HTTPS for Paperless-ngx connections
   - Protect API tokens (use environment variables, not hardcoded values)
   - Run as non-root user (Docker images already configured)
   - Use `MCP_AUTH_TOKEN` when exposing HTTP transport

2. **Performance**:
   - Configure appropriate resource limits in docker-compose.yml
   - Monitor memory usage for large document operations
   - Use pagination for list operations

3. **Reliability**:
   - Enable Docker restart policies (`restart: unless-stopped`)
   - Monitor health check endpoints
   - Set up logging aggregation

### Docker Deployment

The provided Dockerfile uses multi-stage builds for:
- Small image size (< 20MB runtime)
- Static binary with no runtime dependencies
- Non-root user (UID 1000)
- Health checks configured

Resource limits in docker-compose.yml:
- CPU: 1 core maximum, 0.25 core reserved
- Memory: 256MB maximum, 128MB reserved

Adjust these based on your workload.

### Monitoring and Logging

- **Structured Logging**: All logs use structured format (slog)
- **Log Levels**: Configure via `LOG_LEVEL` environment variable
- **Health Checks**: Available at `/health` endpoint (HTTP mode only)
- **Metrics**: Check Docker container stats: `docker stats paperless-mcp-server`

## Development

### Project Structure

```
paperless-mcp-go/
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── mcp/             # MCP server implementation
│   │   ├── server.go    # Server setup and registration
│   │   ├── tools.go     # Tool registration
│   │   ├── *_handlers.go # Tool handler implementations
│   │   └── transport.go # Transport layer
│   └── paperless/       # Paperless API client
│       ├── client.go    # HTTP client and API methods
│       ├── types.go     # Data type definitions
│       └── errors.go    # Error handling
├── Dockerfile
├── docker-compose.yml
└── README.md
```

### Adding New Tools

1. **Add API Client Method** (`internal/paperless/client.go`):
   ```go
   func (c *Client) YourMethod(ctx context.Context, ...) (*Type, error) {
       // Implementation
   }
   ```

2. **Create Handler** (appropriate `*_handlers.go` file):
   ```go
   func (s *Server) handleYourTool(ctx context.Context, args map[string]interface{}) (interface{}, error) {
       // Implementation
   }
   ```

3. **Register Tool** (`internal/mcp/tools.go`):
   ```go
   err = s.RegisterTool(Tool{
       Name:        "your_tool",
       Description: "Description of your tool",
       InputSchema: map[string]interface{}{...},
       Handler:     s.handleYourTool,
   })
   ```

### Code Style Guidelines

- **Constants**: Use descriptive constant names for all magic values
- **Logging**: Use structured logging with slog (Debug/Info/Error levels)
- **Error Handling**: Wrap errors with `fmt.Errorf("message: %w", err)`
- **Validation**: Validate all inputs in handlers before calling API methods
- **Pagination**: Support pagination for all list operations

## Troubleshooting

### Common Issues

**Issue**: `Failed to load configuration: environment variable PAPERLESS_URL is required`
- **Solution**: Ensure all required environment variables are set. Check your `.env` file.

**Issue**: `Request failed: 401 Unauthorized`
- **Solution**: Verify your `PAPERLESS_TOKEN` is correct and has appropriate permissions.

**Issue**: `Connection refused` when using HTTP transport
- **Solution**: Ensure `MCP_HTTP_PORT` is not already in use and firewall allows connections.

### Debug Logging

Enable debug logging for detailed information:

```bash
export LOG_LEVEL=debug
./paperless-mcp
```

Debug logs include:
- API request/response details
- Tool invocation parameters
- Detailed error information

### Health Check

For HTTP transport mode, check server health:

```bash
curl http://localhost:8080/health
```

Should return health status and server information.

## License

See [LICENSE](LICENSE) file for details.

## Links

- [Paperless-ngx Documentation](https://docs.paperless-ngx.com/)
- [Model Context Protocol](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/mark3labs/mcp-go)

---

**Note**: This is an unofficial community project and is not affiliated with or endorsed by the Paperless-ngx project.
