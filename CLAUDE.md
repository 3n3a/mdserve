# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`mdserve` is a minimal Python HTTP server that renders `.md` files as HTML using Pico CSS. It serves markdown files from the **parent directory** (`CONTENT_DIR = Path(__file__).resolve().parent.parent`) — so this repo itself must be placed inside the folder of `.md` files you want to browse.

## Commands

This project uses [uv](https://docs.astral.sh/uv/) for dependency management (Python 3.14 required).

```bash
# Install dependencies
uv sync

# Run the server (default port 8000, default content dir is parent of script)
uv run server.py

# Run with a custom content directory and/or port
uv run server.py /path/to/notes --port 9000

# Install globally as a CLI tool
uv tool install .

# Then run from anywhere
mdserve /path/to/notes --port 9000
```

There are no tests and no linter configured.

## Architecture

All logic lives in a single file: `server.py`.

- `CONTENT_DIR` is a module-level global set at startup by `main()` — it defaults to one level above the repo root but can be overridden via the `content_dir` positional CLI argument.
- `collect_md_files()` recursively walks `CONTENT_DIR`, skipping hidden paths, and groups files by subfolder.
- `extract_title()` reads the first `# heading` from a file; falls back to the filename stem.
- `Handler.do_GET()` handles two cases: `/` renders an index page (`render_index()`); any other path is resolved against `CONTENT_DIR`, converted from Markdown to HTML via `markdown.markdown()`, and wrapped in `render_page()`.
- The HTML shell in `render_page()` injects an inline sidebar (`render_sidebar()`) and loads Pico CSS from a CDN.
- Markdown extensions enabled: `tables`, `fenced_code`, `toc`, `sane_lists`, `smarty`.
