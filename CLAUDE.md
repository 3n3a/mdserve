# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`mdserve` is a minimal Python HTTP server (FastAPI + uvicorn) that renders `.md` files as HTML using Pico CSS. It serves markdown files from the **current working directory** (or a specified directory) when invoked.

## Commands

This project uses [uv](https://docs.astral.sh/uv/) for dependency management (Python 3.14 required).

```bash
# Install dependencies
uv sync

# Run the server (serves .md files from the current directory, default port 8000)
uv run -m mdserve.server

# Run with a custom port
uv run -m mdserve.server --port 9000

# Serve a specific directory
uv run -m mdserve.server ~/notes --port 9000

# Install globally as a CLI tool
uv tool install .

# Then run from anywhere
mdserve [content_dir] --port 9000
```

There are no tests and no linter configured.

## Architecture

All logic lives in a single file: `server.py`.

- Built on FastAPI + uvicorn. All logic lives in `src/mdserve/server.py`.
- `CONTENT_DIR` is a module-level global defaulting to `Path.cwd()`, overridable via CLI positional arg.
- `collect_md_files()` recursively walks `CONTENT_DIR`, skipping hidden paths, and groups files by subfolder.
- `extract_title()` reads the first `# heading` from a file; falls back to the filename stem.
- Two FastAPI routes: `/` renders an index page (`render_index()`); `/{file_path:path}` resolves against `CONTENT_DIR`, converts from Markdown to HTML via `markdown.markdown()`, and wraps in `render_page()`.
- The HTML shell in `render_page()` injects an inline sidebar (`render_sidebar()`) and loads Pico CSS from a CDN.
- File reading tries UTF-8-sig, UTF-8, then Latin-1 fallback.
- Markdown extensions enabled: `tables`, `fenced_code`, `toc`, `sane_lists`, `smarty`.
