package changelog

import (
	"context"
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

func (clg *ChangelogGenerator) Generate(ctx context.Context) (changelog *Changelog, err error) {
	slog.Info("Generating changelog", "repositories", clg.config.Repositories, "organization", clg.config.Organization)
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
				} `graphql:"releases(first: 20)"`
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
			slog.Debug("Release", "tag", release.Tag.Name, "date", release.CreatedAt, "name", release.Name, "description", release.Description)
			for _, pr := range release.Tag.AssociatedPullRequests.Nodes {
				slog.Debug("PR", "number", pr.Number, "author", pr.Author.Login)
			}
		}
	}

	return &Changelog{
		Releases: []Release{
			{
				Tag:   "1.1.0",
				Date:  "2020-01-02",
				Title: "Added new feature",
			},
			{
				Tag:   "1.0.0",
				Date:  "2020-01-01",
				Title: "Initial release",
			},
		},
	}, nil
}
