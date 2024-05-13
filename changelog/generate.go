package changelog

import (
	"context"
	"github.com/prodyna/changelog-json/changelog/output"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"log/slog"
	"strings"
)

type Config struct {
	GitHubToken  string
	Repositories string
	Organization string
}

type Tag struct {
	Name        string
	Description string
	Repository  string
}

type ChangelogGenerator struct {
	config *Config
}

func New(config Config) (*ChangelogGenerator, error) {
	c := &ChangelogGenerator{
		config: &config,
	}
	slog.Info("Changlog Generator", "organization", config.Organization, "repositories", config.Repositories)
	return c, nil
}

func (clg *ChangelogGenerator) Generate(ctx context.Context) (changelog *output.Changelog, err error) {
	slog.Info("Generating changelog", "repositories", clg.config.Repositories, "organization", clg.config.Organization)
	changelog = &output.Changelog{
		Releases: &[]output.Release{},
	}
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: clg.config.GitHubToken},
	)
	httpClient := oauth2.NewClient(ctx, src)
	client := githubv4.NewClient(httpClient)

	var query struct {
		Organization struct {
			Repository struct {
				Name     string
				Releases struct {
					Nodes []struct {
						Name        string
						CreatedAt   string
						Description string
						Tag         struct {
							Name                   string
							AssociatedPullRequests struct {
								Nodes []struct {
									Number int
									Author struct {
										Login string
									}
								}
							} `graphql:"associatedPullRequests(first: 99)"`
						}
					}
				} `graphql:"releases(first: 10)"`
			} `graphql:"repository(name: $repository)"`
		} `graphql:"organization(login: $organization)"`
	}

	repositories := strings.Split(clg.config.Repositories, ",")
	for _, repository := range repositories {
		slog.Debug("Generating changelog for repository", "repository", repository)
		variables := map[string]interface{}{
			"organization": githubv4.String(clg.config.Organization),
			"repository":   githubv4.String(repository),
		}
		err = client.Query(ctx, &query, variables)
		if err != nil {
			slog.Error("unable to query github", "error", err, "repository", repository)
			return nil, err
		}

		for _, release := range query.Organization.Repository.Releases.Nodes {
			slog.Debug("Release", "tag", release.Tag.Name, "date", release.CreatedAt, "name", release.Name, "description.len", len(release.Description))
			entry := output.Entry{
				Tag:         release.Tag.Name,
				Name:        release.Name,
				Component:   repository,
				Description: release.Description,
			}
			slog.Info("Adding entry", "tag", entry.Tag, "name", entry.Name, "component", entry.Component, "description.len", len(entry.Description))
			changelog.AddEntry(entry)
		}
	}

	return changelog, nil
}
