package main

import (
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/htmx/router"
	"log/slog"
	"os"
)

func main() {
	cmd.InitFlags()
	// Run Comigo server.
	if err := router.RunServer(); err != nil {
		slog.Error("Failed to start server!", "details", err.Error())
		os.Exit(1)
	}
}
