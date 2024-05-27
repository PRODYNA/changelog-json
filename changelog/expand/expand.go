package expand

import (
	"log/slog"
	"regexp"
)

func ExpandLinks(description string) string {
	slog.Debug("Expanding links")
	r := regexp.MustCompile(`(https://[^ \r\n]+)`)
	return r.ReplaceAllString(description, "[$1]($1)")
}
