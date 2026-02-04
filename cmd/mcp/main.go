package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/memohai/memoh/internal/logger"
	"github.com/memohai/memoh/internal/mcp"
	"github.com/memohai/memoh/internal/version"
	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := gomcp.NewServer(
		&gomcp.Implementation{Name: "memoh-mcp", Version: version.GetInfo()},
		nil,
	)
	mcp.RegisterTools(server)
	err := server.Run(ctx, &gomcp.StdioTransport{})
	if ctx.Err() != nil {
		return
	}
	if err == nil {
		logger.Warn("mcp server exited without error; waiting for shutdown signal")
		<-ctx.Done()
		return
	}
	if errors.Is(err, io.EOF) {
		logger.Warn("mcp stdio closed; waiting for shutdown signal")
		<-ctx.Done()
		return
	}
	logger.Error("mcp server failed", slog.Any("error", err))
	os.Exit(1)
}
