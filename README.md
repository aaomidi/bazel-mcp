# Bazel MCP Server

This project provides a Model Context Protocol (MCP) server that exposes tools for interacting with Bazel projects.

**Status:** 🚧 Under Active Development 🚧

Features and interfaces may change.

## Features

*   **Reverse Dependencies (`reverse-dependencies`):** Find all Bazel targets that depend on a given target or file path within a project.

## Usage with Claude Desktop

This server is designed to be used as a tool provider for applications like the Claude Desktop app.

### 1. Installation

Install the server using `go install`:

```bash
# Ensure your Go environment is set up (GOPATH, GOBIN, etc.)
go install github.com/aaomidi/mcp-bazel@latest
```

This command compiles and installs the `mcp-bazel` binary. By default, it's placed in your `$GOPATH/bin` directory or `$HOME/go/bin` if `GOPATH` is not set. Make sure this location is in your system's `PATH` or note the full path to the binary for the next step.

### 2. Configuration

Add the following configuration to your Claude Desktop settings file (`claude_desktop_config.json`). You can usually find this file via the app's settings menu ("Open Config Folder").

Modify the `"command"` path to point to the actual location where `go install` placed the `mcp-bazel` binary on your system.

```json
{
  "mcpServers": {
    "mcp-bazel": {
      "command": "/path/to/your/go/bin/mcp-bazel"
    }
  }
}
```

### 3. Restart Claude Desktop

After saving the configuration file, restart the Claude Desktop app. The Bazel tools should now be available.
