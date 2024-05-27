package expand

import "testing"

func TestExpandLinks(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        string
	}{
		{
			name:        "no links",
			description: "no links",
			want:        "no links",
		},
		{
			name:        "pull request",
			description: "this is a pull request https://github.com/PRODYNA-YASM/yasm-backend/pull/549",
			want:        "this is a pull request [**#PR549**](https://github.com/PRODYNA-YASM/yasm-backend/pull/549)",
		},
		{
			name:        "changelog",
			description: "this is a changelog https://github.com/PRODYNA-YASM/yasm-backend/compare/1.16.4...1.19.0",
			want:        "this is a changelog [**#1.16.4...1.19.0**](https://github.com/PRODYNA-YASM/yasm-backend/compare/1.16.4...1.19.0)",
		},
		{
			name:        "github user",
			description: "this is a github user @dkrizic",
			want:        "this is a github user [**@dkrizic**](https://github.com/dkrizic)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExpandLinks(tt.description); got != tt.want {
				t.Errorf("ExpandLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}
