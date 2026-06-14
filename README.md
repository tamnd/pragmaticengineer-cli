# pragmaticengineer

Browse [The Pragmatic Engineer](https://newsletter.pragmaticengineer.com/) newsletter from the command line.

`pragmaticengineer` is a single pure-Go binary. It reads public data from The Pragmatic Engineer
Substack over plain HTTPS, shapes it into clean records, and prints output that pipes
into the rest of your tools. No API key, nothing to run alongside it.

## Install

```bash
go install github.com/tamnd/pragmaticengineer-cli/cmd/pragmaticengineer@latest
```

Or grab a prebuilt binary from the [releases](https://github.com/tamnd/pragmaticengineer-cli/releases), or run
the container image:

```bash
docker run --rm ghcr.io/tamnd/pragmaticengineer:latest --help
```

## Usage

```bash
# List the 25 most recent posts (default)
pragmaticengineer top

# Table output in the terminal
pragmaticengineer top -o table

# Get just the URLs
pragmaticengineer top -o url

# Fetch 10 posts starting at offset 25 (second page)
pragmaticengineer top --limit 10 --offset 25

# As JSON for jq
pragmaticengineer top -o json | jq '.[] | select(.audience == "free")'

# CSV for spreadsheets
pragmaticengineer top -o csv > posts.csv
```

Every command shares one output contract: `-o table|json|jsonl|csv|tsv|url|raw`,
`--fields` to pick columns, `--template` for a custom line, and `-n` to limit.
The default adapts to where output goes (a table on a terminal, JSONL in a
pipe), so the same command reads well by hand and parses cleanly downstream.

## Commands

| Command | Description |
|---------|-------------|
| `top` | List recent posts from The Pragmatic Engineer |
| `version` | Show version information |

## Flags for `top`

```
--limit int     number of posts to fetch (default 25)
--offset int    offset for pagination (default 0)
```

## Global flags

```
-o, --output string    output format: table|json|jsonl|csv|tsv|url|raw (default "auto")
-n, --limit int        limit number of records (0 = all)
    --fields strings   comma-separated columns to include
    --no-header        omit header row
    --template string  Go text/template per record
    --timeout duration per-request timeout (default 30s)
    --delay duration   minimum spacing between requests
    --retries int      retry attempts on 429/5xx (default 3)
```

## Post fields

| Field | Description |
|-------|-------------|
| `rank` | Position in the listing (1-based) |
| `date` | Publication date (YYYY-MM-DD) |
| `audience` | `free` or `paid` |
| `title` | Post title |
| `subtitle` | Post subtitle / summary |
| `url` | Canonical URL |

## Serve it

The same operations are available over HTTP and as an MCP tool set for agents,
with no extra code:

```bash
pragmaticengineer serve --addr :7777    # GET /v1/top returns NDJSON
pragmaticengineer mcp                   # speak MCP over stdio
```

## Development

```
cmd/pragmaticengineer/   thin main: hands cli.NewApp to kit.Run
cli/                     assembles the kit App from the pragmaticengineer domain
pragmaticengineer/       the library: HTTP client, data models, and domain.go (the driver)
docs/                    tago documentation site
```

```bash
make build      # ./bin/pragmaticengineer
make test       # go test ./...
make vet        # go vet ./...
```

## Releasing

Push a version tag and GitHub Actions runs GoReleaser, which builds the
archives, Linux packages, the multi-arch GHCR image, checksums, SBOMs, and a
cosign signature:

```bash
git tag v0.1.0
git push --tags
```

## License

Apache-2.0. See [LICENSE](LICENSE).
