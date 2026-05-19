# mdserve

A minimal Go HTTP server that renders Markdown files as HTML with a sidebar navigation, [Pico CSS](https://picocss.com/), and [Mermaid](https://mermaid.js.org/) diagrams.

Single static binary — no runtime dependencies.

## Install

Download the latest binary for your platform from the [releases page](https://github.com/3n3a/mdserve/releases), or build from source:

```bash
go install github.com/3n3a/mdserve/cmd/mdserve@latest
```

## Usage

```bash
mdserve [flags] [content_dir]
```

| Flag | Default | Description |
|---|---|---|
| `content_dir` (positional) | `.` | Directory to scan for `.md` files |
| `--host` | `127.0.0.1` | Host to bind to. Use `0.0.0.0` to expose on the network. |
| `--port` | `8000` | Port to listen on (also reads `$PORT`) |
| `--user` | _(empty)_ | Basic-auth username. Enables auth when set. Also reads `$MDSERVE_USER`. |
| `--pass` | _(empty)_ | Basic-auth password. Also reads `$MDSERVE_PASS`. |
| `--offline` | `false` | Serve Pico CSS and Mermaid from the binary instead of the jsDelivr CDN. Also reads `$MDSERVE_OFFLINE`. |
| `--version` | — | Print version and exit |

**Examples:**

```bash
# Serve current directory
mdserve

# Serve a specific folder on a custom port
mdserve ~/notes --port 9000

# Expose publicly with basic auth (see warning below)
mdserve --host 0.0.0.0 --port 9000 --user me --pass hunter2

# Air-gapped / no outbound calls from the browser
mdserve --offline
```

## Public exposure

You can bind to a public interface and protect access with `--user`/`--pass`. This adds HTTP Basic Auth; over plain HTTP credentials are sent base64-encoded only — put it behind a TLS-terminating reverse proxy if you actually need security.

**No warranty.** This tool was not designed as a hardened public service. Run it on the open internet at your own risk.

## Build from source

Requires Go 1.24+.

```bash
# Fetch bundled CSS/JS (needed only if you want --offline to work)
bash scripts/fetch-assets.sh        # or: pwsh scripts/fetch-assets.ps1

go build ./cmd/mdserve
./mdserve
```

## Releases

Pushing a `v*` tag triggers a GitHub Actions workflow that cross-compiles binaries for:

- Linux amd64, Linux arm64
- macOS amd64, macOS arm64
- Windows amd64

and uploads them to the GitHub release for that tag.
