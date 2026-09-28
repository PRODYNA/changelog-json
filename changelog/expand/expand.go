package expand

import (
	"log/slog"
	"regexp"
	"strings"
)

const DefaultServerUrl = "https://github.com"

// ExpandLinks expands links, pull request references and user handles in the
// description. serverUrl is the base URL of the GitHub installation, for
// example https://github.com or https://github.example.com for GitHub Enterprise.
func ExpandLinks(description string, serverUrl string) string {
	slog.Debug("Expanding links", "serverUrl", serverUrl)

	server := strings.TrimSuffix(strings.TrimSpace(serverUrl), "/")
	if server == "" {
		server = DefaultServerUrl
	}
	quoted := regexp.QuoteMeta(server)

	// <server>/PRODYNA-YASM/yasm-backend/pull/549 -> [**#PR549**](<server>/PRODYNA-YASM/yasm-backend/pull/549)
	r := regexp.MustCompile(quoted + "/(.*?)/pull/(\\d+)")
	description = r.ReplaceAllString(description, "[**#PR$2**]("+server+"/$1/pull/$2)")

	// <server>/PRODYNA-YASM/yasm-backend/compare/1.16.4...1.19.0 -> [**#1.16.4...1.19.0**](<server>/PRODYNA-YASM/yasm-backend/compare/1.16.4...1.19.0)
	r = regexp.MustCompile(quoted + "/(.*?)/compare/(.*)")
	description = r.ReplaceAllString(description, "[**#$2**]("+server+"/$1/compare/$2)")

	// @dkrizic -> [**@dkrizic**](<server>/dkrizic)
	// regex that matches github usernames including dashes
	r = regexp.MustCompile("@([a-zA-Z0-9-]+)")
	description = r.ReplaceAllString(description, "[**@$1**]("+server+"/$1)")

	// <!-- blabla --> -> ""
	r = regexp.MustCompile("<!--.*?-->")
	description = r.ReplaceAllString(description, "")

	return description
}
