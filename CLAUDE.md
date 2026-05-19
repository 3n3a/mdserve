# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`mdserve` is a small Go HTTP server that walks a directory, renders `.md` files as HTML with a sidebar navigation, and ships as a single static binary. Pico CSS and Mermaid are loaded from a CDN by default or embedded into the binary via `--offline`.

## Commands

Requires Go 1.24+.

```bash
# Fetch bundled CSS/JS into internal/assets/ (only needed for --offline; CI runs this automatically)
bash scripts/fetch-assets.sh        # POSIX
pwsh scripts/fetch-assets.ps1       # Windows

# Build & run
go build ./cmd/mdserve
./mdserve [content_dir] [--port PORT] [--host HOST] [--user U --pass P] [--offline]

# Vet & test
go vet ./...
go test ./...

# Cross-compile (example)
GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o dist/mdserve-linux-arm64 ./cmd/mdserve
```

## Commit convention

Use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short description>
```

- **type:** `feat`, `fix`, `chore`, `refactor`, `docs`, `test`, `build`, `ci`, `perf`, `style`
- **scope:** the area touched (e.g. `markdown`, `server`, `collector`, `ci`, `release`, `structure`)
- **description:** imperative mood, lowercase, no trailing period

Examples from this repo's history:

```
feat(mdserve): implement mdserve
feat(fastapi): rewrite to use fastapi
fix(markdown): markdown list not correctly rendered
chore(structure): separate code by concern
```

## Architecture

Code is separated by concern under standard Go layout:

```
cmd/mdserve/         CLI entrypoint, flag parsing, graceful shutdown
internal/collector/  Walk dir, group .md files by folder, extract H1 titles
internal/fileio/     Multi-encoding read (UTF-8 BOM strip, UTF-8 validate, Latin-1 fallback)
internal/markdown/   goldmark (GFM + Typographer) + mermaid post-process
internal/render/     html/template page shell, sidebar, index TOC
internal/assets/     embed.FS for pico.min.css + mermaid.esm.min.mjs (gitignored, downloaded at build time)
internal/server/     chi router, handlers, basic-auth middleware
scripts/             fetch-assets.{sh,ps1} downloads pinned CDN versions
.github/workflows/   ci.yml (vet/build/test) and release.yml (5-platform cross-build on tag push)
```

Key behavior:

- `GET /` renders an index TOC. `GET /*` resolves the path against the content dir (with `filepath.Rel` traversal check), converts via goldmark, and wraps in the shared page shell.
- Hidden files/dirs (`.git`, `.github`, etc.) are skipped in the walk.
- Goldmark renders fenced ` ```mermaid ` blocks as `<pre><code class="language-mermaid">…`; a regex post-pass rewrites them to `<pre class="mermaid">` so Mermaid.js can take over.
- Asset URLs are computed once per request: CDN by default, `/_assets/…` (served from `embed.FS`) when `--offline` is set. CI fetches the real CSS/JS before building so released binaries always have them; the committed placeholders are zero-byte and trigger a startup warning if `--offline` is used without running the fetch script.
- Auth (`--user` + `--pass`) is a login page + session cookie. `chi` middleware in `internal/server/auth.go` checks for a `mdserve_session` cookie (HMAC-SHA256 over an 8-byte expiry, signed with a per-process random secret). Unauthenticated requests redirect to `/login?next=...`. Public routes (no auth): `/login`, `/logout`, `/_assets/*`. Credential comparison uses `subtle.ConstantTimeCompare`.
