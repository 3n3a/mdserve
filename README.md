# mdserve

A minimal local web server that renders Markdown files as HTML with a sidebar navigation. Built with Python's stdlib HTTP server and [Pico CSS](https://picocss.com/).

## Installation

Requires [uv](https://docs.astral.sh/uv/).

```bash
uv version --bump minor
uv tool install .
```

## Usage

```bash
mdserve [content_dir] [--port PORT]
```

| Argument | Default | Description |
|---|---|---|
| `content_dir` | parent of script | Directory to scan for `.md` files |
| `--port` | `8000` | Port to listen on (also reads `$PORT`) |

**Examples:**

```bash
# Serve current directory
mdserve .

# Serve a specific folder on a custom port
mdserve ~/notes --port 9000
```

Then open `http://127.0.0.1:8000` in your browser.

## Without installing

```bash
uv run -m mdserve.server [content_dir] [--port PORT]
```
