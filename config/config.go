package config

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	keyVerbose                 = "verbose"
	keyVerboseEnvironment      = "VERBOSE"
	keyGithubToken             = "github-token"
	keyGithubTokenEnvironment  = "GITHUB_TOKEN"
	keyRepositories            = "repositories"
	keyRepositoriesEnvironment = "REPOSITORIES"
	keyOrganization            = "organization"
	keyOrganizationEnvironment = "ORGANIZATION"
	keyOutputFile              = "output-file"
	keyOutputFileEnvironment   = "OUTPUT_FILE"
	keyExpandLinks             = "expand-links"
	keyExpandLinksEnvironment  = "EXPAND_LINKS"

	keyGithubGraphQLUrl            = "github-graphql-url"
	keyGithubGraphQLUrlEnvironment = "GITHUB_GRAPHQL_URL"
	keyGithubServerUrl             = "github-server-url"
	keyGithubServerUrlEnvironment  = "GITHUB_SERVER_URL"

	defaultGithubGraphQLUrl = "https://api.github.com/graphql"
	defaultGithubServerUrl  = "https://github.com"
)

type Config struct {
	Verbose          *int
	GithubToken      string
	GithubGraphQLUrl string
	GithubServerUrl  string
	Repositories     string
	Organization     string
	OutputFile       string
	ExpandLinks      bool
}

func New() (*Config, error) {
	c := Config{}
	verbose := flag.Int(keyVerbose, lookupEnvOrInt(keyVerboseEnvironment, 0), "Verbosity level, 0=info, 1=debug. Overrides the environment variable VERBOSE.")
	flag.StringVar(&c.GithubToken, keyGithubToken, lookupEnvOrString(keyGithubTokenEnvironment, ""), "The GitHub Token to use for authentication.")
	flag.StringVar(&c.GithubGraphQLUrl, keyGithubGraphQLUrl, lookupEnvOrString(keyGithubGraphQLUrlEnvironment, defaultGithubGraphQLUrl), "The GitHub GraphQL API endpoint, for example https://github.example.com/api/graphql for GitHub Enterprise.")
	flag.StringVar(&c.GithubServerUrl, keyGithubServerUrl, lookupEnvOrString(keyGithubServerUrlEnvironment, defaultGithubServerUrl), "The GitHub server URL used to expand links, for example https://github.example.com for GitHub Enterprise.")
	flag.StringVar(&c.Repositories, keyRepositories, lookupEnvOrString(keyRepositoriesEnvironment, ""), "The repositories to generate changelog for.")
	flag.StringVar(&c.Organization, keyOrganization, lookupEnvOrString(keyOrganizationEnvironment, ""), "The organization to generate changelog for.")
	flag.StringVar(&c.OutputFile, keyOutputFile, lookupEnvOrString(keyOutputFileEnvironment, ""), "The output file to write the changelog to.")
	flag.BoolVar(&c.ExpandLinks, keyExpandLinks, lookupEnvOrBool(keyExpandLinksEnvironment, "true"), "Expand links in the changelog.")
	flag.Parse()

	level := slog.LevelError
	switch *verbose {
	case 0:
		level = slog.LevelError
	case 1:
		level = slog.LevelWarn
	case 2:
		level = slog.LevelInfo
	case 3:
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})))

	if c.GithubToken == "" {
		return nil, fmt.Errorf("missing required environment variable: %s", keyGithubToken)
	}

	if c.Repositories == "" {
		return nil, fmt.Errorf("missing required environment variable: %s", keyRepositories)
	}

	if c.Organization == "" {
		return nil, fmt.Errorf("missing required environment variable: %s", keyOrganization)
	}

	if c.OutputFile == "" {
		return nil, fmt.Errorf("missing required environment variable: %s", keyOutputFile)
	}

	graphQLUrl, err := normalizeGraphQLUrl(c.GithubGraphQLUrl)
	if err != nil {
		return nil, err
	}
	c.GithubGraphQLUrl = graphQLUrl

	serverUrl, err := normalizeServerUrl(c.GithubServerUrl)
	if err != nil {
		return nil, err
	}
	c.GithubServerUrl = serverUrl

	slog.Debug("GitHub endpoints", "graphql", c.GithubGraphQLUrl, "server", c.GithubServerUrl)

	return &c, nil
}

// normalizeGraphQLUrl accepts a GraphQL endpoint, a REST API base URL or the
// plain host URL of a GitHub (Enterprise) installation and returns the
// GraphQL endpoint to use.
func normalizeGraphQLUrl(raw string) (string, error) {
	if raw == "" {
		return defaultGithubGraphQLUrl, nil
	}
	u, err := parseUrl(raw, keyGithubGraphQLUrl)
	if err != nil {
		return "", err
	}
	path := strings.TrimSuffix(u.Path, "/")
	switch {
	case strings.HasSuffix(path, "/graphql"):
		// already a GraphQL endpoint
	case u.Host == "api.github.com":
		path = path + "/graphql"
	case strings.HasSuffix(path, "/api/v3"):
		path = strings.TrimSuffix(path, "/v3") + "/graphql"
	case strings.HasSuffix(path, "/api"):
		path = path + "/graphql"
	default:
		path = path + "/api/graphql"
	}
	u.Path = path
	return u.String(), nil
}

// normalizeServerUrl returns the scheme and host of the GitHub (Enterprise)
// web frontend, without a trailing slash.
func normalizeServerUrl(raw string) (string, error) {
	if raw == "" {
		return defaultGithubServerUrl, nil
	}
	u, err := parseUrl(raw, keyGithubServerUrl)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(u.Scheme+"://"+u.Host+u.Path, "/"), nil
}

func parseUrl(raw string, key string) (*url.URL, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return nil, fmt.Errorf("invalid value for %s: %w", key, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid value for %s: %q is not an absolute URL", key, raw)
	}
	return u, nil
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func lookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}

func lookupEnvOrBool(key string, defaultVal string) bool {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatalf("LookupEnvOrBool[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal == "true"
}
