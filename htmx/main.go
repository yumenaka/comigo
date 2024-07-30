package main

import (
	"github.com/yumenaka/comi/htmx/router"
	"log/slog"
	"os"
)

func main() {
	// Run your server.
	if err := router.RunServer(); err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}
}
