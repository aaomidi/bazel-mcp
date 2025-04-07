package main

import (
	"log"

	"github.com/aaomidi/mcp-bazel/tools/rdeps"

	"github.com/mark3labs/mcp-go/server"
)

func _main() error {
	s := server.NewMCPServer(
		"bazel-mcp",
		"0.0.1",
	)

	// Register tools
	s.AddTool(rdeps.Tool, rdeps.Handler)

	return server.ServeStdio(s)
}

func main() {
	if err := _main(); err != nil {
		log.Fatalf("%+v", err)
	}
}
