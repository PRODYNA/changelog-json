package main

import (
	"github.com/prodyna/changelog-json/config"
	"log/slog"
	"os"
)

func main() {
	_, err := config.New()
	if err != nil {
		slog.Error("unable to load config", "error", err)
		os.Exit(1)
	}

	slog.Debug("config loaded")
}
