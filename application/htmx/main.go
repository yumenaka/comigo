package main

import (
	"log/slog"
	"os"
)

func main() {
	// Run your server.
	if err := runServer(); err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}
}
