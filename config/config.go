package config

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
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
)

type Config struct {
	Verbose      *int
	GithubToken  string
	Repositories string
	Organization string
	OutputFile   string
	ExpandLinks  bool
}

func New() (*Config, error) {
	c := Config{}
	verbose := flag.Int(keyVerbose, lookupEnvOrInt(keyVerboseEnvironment, 0), "Verbosity level, 0=info, 1=debug. Overrides the environment variable VERBOSE.")
	flag.StringVar(&c.GithubToken, keyGithubToken, lookupEnvOrString(keyGithubTokenEnvironment, ""), "The GitHub Token to use for authentication.")
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

	return &c, nil
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
