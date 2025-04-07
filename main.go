package main

import (
	"log"

	"github.com/aaomidi/mcp-bazel/tools/build"
	"github.com/aaomidi/mcp-bazel/tools/rdeps"
	"github.com/aaomidi/mcp-bazel/tools/sources"
	"github.com/aaomidi/mcp-bazel/tools/test"

	"github.com/mark3labs/mcp-go/server"
)

func _main() error {
	s := server.NewMCPServer(
		"bazel-mcp",
		"0.0.1",
	)

	// Register tools
	s.AddTool(rdeps.Tool, rdeps.Handler)
	s.AddTool(build.Tool, build.Handler)
	s.AddTool(test.Tool, test.Handler)
	s.AddTool(sources.Tool, sources.Handler)

	return server.ServeStdio(s)
}

func main() {
	if err := _main(); err != nil {
		log.Fatalf("%+v", err)
	}
}
