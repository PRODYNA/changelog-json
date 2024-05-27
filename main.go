package main

import (
	"context"
	"github.com/prodyna/changelog-json/changelog"
	"github.com/prodyna/changelog-json/config"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()

	c, err := config.New()
	if err != nil {
		slog.Error("unable to load config", "error", err)
		os.Exit(1)
	}

	slog.Debug("config loaded")

	generator, err := changelog.New(changelog.Config{
		GitHubToken:  c.GithubToken,
		Repositories: c.Repositories,
		Organization: c.Organization,
		ExpandLinks:  c.ExpandLinks,
	})
	if err != nil {
		slog.Error("unable to create changelog generator", "error", err)
		os.Exit(1)
	}

	cl, err := generator.Generate(ctx)
	if err != nil {
		slog.Error("unable to generate changelog", "error", err)
		os.Exit(1)
	}

	output, err := cl.RenderJSON()
	if err != nil {
		slog.Error("unable to render changelog", "error", err)
		os.Exit(1)
	}

	// write changelog to file
	err = os.WriteFile(c.OutputFile, output, 0644)
	if err != nil {
		slog.Error("unable to write changelog to file", "error", err)
		os.Exit(1)
	}
}
