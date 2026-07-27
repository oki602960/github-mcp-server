# GitHub MCP Server

An official GitHub MCP (Model Context Protocol) server written in Go. Connects AI tools (VS Code, Claude, Cursor, Windsurf, etc.) directly to GitHub's API.

## How to run

```bash
GOTOOLCHAIN=local ./github-mcp-server http --port 5000
```

The workflow "Start application" is configured to run this automatically on port 5000.

## Authentication

The HTTP server does **not** require a token at startup. Each client request must supply a `GITHUB_PERSONAL_ACCESS_TOKEN` via the `Authorization: Bearer <token>` HTTP header.

## Build

```bash
GOTOOLCHAIN=local go build -o github-mcp-server ./cmd/github-mcp-server/
```

Go 1.25 is required (installed as module `go-1.25`). The `go.mod` minimum was pinned to `1.25.5` (the version available on Replit) from the original `1.25.12`.

## CORS

CORS is fully open (`Access-Control-Allow-Origin: *`) in `pkg/http/middleware/cors.go` — no restrictions.

## User preferences

- Run the HTTP server on port 5000 for Replit webview compatibility.
