package config

import "testing"

func TestNormalizeGraphQLUrl(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{name: "empty", raw: "", want: defaultGithubGraphQLUrl},
		{name: "api.github.com", raw: "https://api.github.com", want: "https://api.github.com/graphql"},
		{name: "api.github.com with graphql", raw: "https://api.github.com/graphql", want: "https://api.github.com/graphql"},
		{name: "enterprise host", raw: "https://github.example.com", want: "https://github.example.com/api/graphql"},
		{name: "enterprise host with slash", raw: "https://github.example.com/", want: "https://github.example.com/api/graphql"},
		{name: "enterprise rest v3", raw: "https://github.example.com/api/v3", want: "https://github.example.com/api/graphql"},
		{name: "enterprise api", raw: "https://github.example.com/api", want: "https://github.example.com/api/graphql"},
		{name: "enterprise graphql", raw: "https://github.example.com/api/graphql", want: "https://github.example.com/api/graphql"},
		{name: "not absolute", raw: "github.example.com", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeGraphQLUrl(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Fatalf("normalizeGraphQLUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("normalizeGraphQLUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeServerUrl(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{name: "empty", raw: "", want: defaultGithubServerUrl},
		{name: "github.com", raw: "https://github.com", want: "https://github.com"},
		{name: "trailing slash", raw: "https://github.example.com/", want: "https://github.example.com"},
		{name: "enterprise", raw: "https://github.example.com", want: "https://github.example.com"},
		{name: "not absolute", raw: "github.example.com", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeServerUrl(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Fatalf("normalizeServerUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got != tt.want {
				t.Errorf("normalizeServerUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
