# ChangeLog JSON

GitHub Action (and standalone CLI) that queries the GitHub API for the releases of one or more
repositories and generates a single, aggregated changelog as a JSON file.

It works against github.com as well as against **GitHub Enterprise (GHE)** installations.

## How it works

1. For every configured repository the action queries the GitHub GraphQL API for the latest
   releases (up to 20 per repository).
2. All releases are merged into a single changelog, grouped by tag and sorted by
   [semantic version](https://semver.org/) in descending order (newest first).
3. Optionally, links, pull request references and user handles inside the release description are
   expanded into Markdown links.
4. The result is written to a JSON file.

## Usage

```yaml
name: Generate changelog

on:
  workflow_dispatch:
  schedule:
    - cron: '0 6 * * *'

jobs:
  changelog:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Generate changelog
        uses: PRODYNA/changelog-json@v1.5
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}
          organization: PRODYNA
          repositories: changelog-json,another-repository
          output-file: CHANGELOG.json
          expand-links: 'true'
```

The token needs read access to the repositories you want to query. For private repositories in
other organizations, the default `GITHUB_TOKEN` is not sufficient — use a Personal Access Token
(PAT) or a GitHub App token with `contents: read` on those repositories.

## Inputs

| Input                | Environment variable  | Required | Default                    | Description                                                                                                                                       |
|----------------------|-----------------------|----------|----------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| `github-token`       | `GITHUB_TOKEN`        | yes      | –                          | Token used to authenticate against the GitHub API.                                                                                                |
| `organization`       | `ORGANIZATION`        | yes      | –                          | The organization (or user) that owns the repositories.                                                                                            |
| `repositories`       | `REPOSITORIES`        | yes      | –                          | Comma-separated list of repository names to query, for example `repo-a,repo-b`.                                                                   |
| `output-file`        | `OUTPUT_FILE`         | yes      | `CHANGELOG.json`           | Path of the JSON file to write.                                                                                                                   |
| `expand-links`       | `EXPAND_LINKS`        | no       | `true`                     | Expand links, pull request URLs and `@user` handles in the release description into Markdown links.                                               |
| `github-graphql-url` | `GITHUB_GRAPHQL_URL`  | no       | `${{ github.graphql_url }}` | GraphQL API endpoint. Defaults to the endpoint of the current GitHub instance, so GHE works out of the box. See [GitHub Enterprise](#github-enterprise). |
| `github-server-url`  | `GITHUB_SERVER_URL`   | no       | `${{ github.server_url }}` | Web URL of the GitHub instance, used when expanding links. Defaults to the URL of the current GitHub instance.                                     |

Additionally, the CLI supports a verbosity setting:

| Flag        | Environment variable | Default | Description                                            |
|-------------|----------------------|---------|--------------------------------------------------------|
| `--verbose` | `VERBOSE`            | `0`     | `0` = error, `1` = warn, `2` = info, `3` = debug.      |

Every input can also be passed as a command line flag when running the binary directly, for
example `--github-graphql-url=https://github.example.com/api/graphql`.

## GitHub Enterprise

By default the action talks to the GitHub instance it is running on, because `github-graphql-url`
and `github-server-url` default to the `github.graphql_url` and `github.server_url` context
values. **When you run this action inside GitHub Enterprise, no extra configuration is needed.**

You only need to set the values explicitly if you want to query a *different* instance than the one
the workflow runs on — for example a workflow running on github.com that reads releases from a GHE
installation:

```yaml
      - name: Generate changelog from GitHub Enterprise
        uses: PRODYNA/changelog-json@v1.5
        with:
          github-token: ${{ secrets.GHE_TOKEN }}
          organization: my-org
          repositories: repo-a,repo-b
          output-file: CHANGELOG.json
          github-graphql-url: https://github.example.com/api/graphql
          github-server-url: https://github.example.com
```

Notes:

- `github-graphql-url` is normalized, so all of the following are accepted and resolve to
  `https://github.example.com/api/graphql`:
  `https://github.example.com`, `https://github.example.com/api`,
  `https://github.example.com/api/v3` and `https://github.example.com/api/graphql`.
- `github-server-url` is only used for link expansion. It must point to the web frontend
  (for example `https://github.example.com`), not to the API.
- Both values must be absolute URLs including the scheme, otherwise the action fails with a
  configuration error.

## Output

The generated file contains all releases, newest first, with one entry per repository
("component") that published that release tag:

```json
{
  "releases": [
    {
      "tag": "v1.4.0",
      "components": [
        {
          "name": "changelog-json",
          "title": "v1.4.0",
          "date": "2024-05-28T10:11:12Z",
          "description": "Fixed a bug in [**#PR42**](https://github.com/PRODYNA/changelog-json/pull/42) by [**@dkrizic**](https://github.com/dkrizic)"
        }
      ]
    }
  ]
}
```

### Link expansion

With `expand-links: 'true'` (the default) the release description is rewritten:

| Input                                                     | Output                                                                                     |
|-----------------------------------------------------------|---------------------------------------------------------------------------------------------|
| `<server>/org/repo/pull/549`                              | `[**#PR549**](<server>/org/repo/pull/549)`                                                  |
| `<server>/org/repo/compare/1.16.4...1.19.0`               | `[**#1.16.4...1.19.0**](<server>/org/repo/compare/1.16.4...1.19.0)`                         |
| `@dkrizic`                                                | `[**@dkrizic**](<server>/dkrizic)`                                                          |
| `<!-- comment -->`                                        | removed                                                                                     |

`<server>` is the value of `github-server-url`, so links point to the correct instance on GHE.

## Running locally

The action is published as a container image, but you can also run it directly:

```bash
go run . \
  --github-token="$GITHUB_TOKEN" \
  --organization=PRODYNA \
  --repositories=changelog-json \
  --output-file=CHANGELOG.json \
  --verbose=3
```

Or via Docker:

```bash
docker run --rm \
  -e GITHUB_TOKEN \
  -e ORGANIZATION=PRODYNA \
  -e REPOSITORIES=changelog-json \
  -e OUTPUT_FILE=/out/CHANGELOG.json \
  -v "$PWD:/out" \
  ghcr.io/prodyna/changelog-json:latest
```

## Development

```bash
go build ./...
go vet ./...
go test ./...
```

The container image is built and pushed to `ghcr.io/prodyna/changelog-json` by the
[build workflow](.github/workflows/build.yaml) on every push to `main` and on every `v*` tag.
