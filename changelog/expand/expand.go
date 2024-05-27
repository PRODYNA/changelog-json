package expand

import (
	"log/slog"
	"regexp"
)

func ExpandLinks(description string) string {
	slog.Debug("Expanding links")

	// https://github.com/PRODYNA-YASM/yasm-backend/pull/549 -> [**#PR549**](https://github.com/PRODYNA-YASM/yasm-backend/pull/549)
	r := regexp.MustCompile("https://github.com/(.*?)/pull/(\\d+)")
	description = r.ReplaceAllString(description, "[**#PR$2**](https://github.com/$1/pull/$2)")

	// https://github.com/PRODYNA-YASM/yasm-backend/compare/1.16.4...1.19.0 -> [**#1.16.4...1.19.0**](https://github.com/PRODYNA-YASM/yasm-backend/compare/1.16.4...1.19.0)
	r = regexp.MustCompile("https://github.com/(.*?)/compare/(.*)")
	description = r.ReplaceAllString(description, "[**#$2**](https://github.com/$1/compare/$2)")
	// 	r := regexp.MustCompile(`(https://[^ \r\n]+)`)
	//	return r.ReplaceAllString(description, "[$1]($1)")

	// @dkrizic -> [**@dkrizic**](https://github.com/dkrizic)
	r = regexp.MustCompile("@(\\w+)")
	description = r.ReplaceAllString(description, "[**@$1**](https://github.com/$1)")

	// <!-- blabla --> -> ""
	r = regexp.MustCompile("<!--.*?-->")
	description = r.ReplaceAllString(description, "")

	return description
}
