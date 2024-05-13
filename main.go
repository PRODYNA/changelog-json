package main

import (
	"github.com/prodyna/changelog-json/changelog"
	"github.com/prodyna/changelog-json/config"
	"log/slog"
	"os"
)

func main() {
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
	})
	if err != nil {
		slog.Error("unable to create changelog generator", "error", err)
		os.Exit(1)
	}

	changelog, err := generator.Generate()
	if err != nil {
		slog.Error("unable to generate changelog", "error", err)
		os.Exit(1)
	}

	output, err := changelog.RenderJSON()
	if err != nil {
		slog.Error("unable to render changelog", "error", err)
		os.Exit(1)
	}

	os.Stdout.Write(output)
}
