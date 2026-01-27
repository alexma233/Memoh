package main

import (
	"context"
	"log"

	"github.com/memohai/memoh/internal/mcp"
	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	commitHash = "unknown"
	version    = "unknown"
)

func main() {
	if version == "unknown" {
		version = "v0.0.0-dev+" + commitHash
	}
	server := gomcp.NewServer(
		&gomcp.Implementation{Name: "memoh-mcp", Version: version},
		nil,
	)
	mcp.RegisterTools(server)
	if err := server.Run(context.Background(), &gomcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
