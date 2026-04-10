"""Minimal webserver that converts .md files to HTML on the fly using Pico CSS."""

from __future__ import annotations

import argparse
import os
from pathlib import Path

import markdown
import uvicorn
from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse

CONTENT_DIR: Path = Path.cwd()

PICO_CSS = "https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css"

MD_EXTENSIONS = ["tables", "fenced_code", "toc", "sane_lists", "smarty"]

app = FastAPI()


def collect_md_files(directory: Path) -> dict[str, list[tuple[str, str, str]]]:
    """Return {folder_name: [(filename, title, url_path), ...]} sorted."""
    folders: dict[str, list[tuple[str, str, str]]] = {}
    for md_file in sorted(directory.rglob("*.md")):
        rel = md_file.relative_to(directory)
        # skip hidden dirs / files
        if any(part.startswith(".") for part in rel.parts):
            continue
        folder = str(rel.parent) if rel.parent != Path(".") else ""
        title = extract_title(md_file)
        url_path = "/" + rel.as_posix()
        folders.setdefault(folder, []).append((md_file.name, title, url_path))
    return dict(sorted(folders.items()))


def extract_title(path: Path) -> str:
    """Extract the first markdown heading or fall back to the filename."""
    try:
        with open(path, encoding="utf-8") as f:
            for line in f:
                line = line.strip()
                if line.startswith("# "):
                    return line.lstrip("# ").strip()
    except Exception:
        pass
    return path.stem


def read_file(path: Path) -> str:
    """Read a file trying multiple encodings."""
    for encoding in ("utf-8-sig", "utf-8", "latin-1"):
        try:
            return path.read_text(encoding=encoding)
        except (UnicodeDecodeError, ValueError):
            continue
    return path.read_text(encoding="latin-1", errors="replace")


def render_sidebar(current_path: str) -> str:
    """Build the sidebar navigation HTML."""
    folders = collect_md_files(CONTENT_DIR)
    parts: list[str] = []
    parts.append('<aside class="sidebar">')
    parts.append('<a href="/" class="sidebar-logo"><strong>mdserve</strong></a>')
    for folder, files in folders.items():
        heading = folder if folder else "Files"
        parts.append(f'<div class="sidebar-heading">{heading}</div><ul>')
        for _fname, title, url_path in files:
            active = ' class="active"' if url_path == current_path else ""
            parts.append(
                f'<li><a href="{url_path}"{active}>{title}</a></li>'
            )
        parts.append("</ul>")
    parts.append("</aside>")
    return "\n".join(parts)


STYLE = """
    *, *::before, *::after { box-sizing: border-box; }
    body {
      display: flex;
      min-height: 100vh;
      margin: 0;
    }
    .sidebar {
      width: 260px;
      flex-shrink: 0;
      padding: 1.25rem 0.75rem;
      background: var(--pico-card-background-color);
      border-right: 1px solid var(--pico-muted-border-color);
      overflow-y: auto;
      position: sticky;
      top: 0;
      height: 100vh;
    }
    .sidebar-logo {
      display: block;
      margin-bottom: 1.25rem;
      padding: 0 0.5rem;
      font-size: 1.1rem;
      text-decoration: none;
      color: var(--pico-color);
    }
    .sidebar-logo:hover {
      color: var(--pico-primary);
    }
    .sidebar-heading {
      font-weight: 600;
      font-size: 0.7rem;
      text-transform: uppercase;
      letter-spacing: 0.06em;
      color: var(--pico-muted-color);
      padding: 0.75rem 0.5rem 0.25rem;
      margin: 0;
    }
    .sidebar ul {
      list-style: none;
      padding: 0;
      margin: 0;
    }
    .sidebar li {
      margin: 0;
      padding: 0;
    }
    .sidebar li a {
      display: block;
      padding: 0.3rem 0.5rem;
      border-radius: 6px;
      text-decoration: none;
      color: var(--pico-color);
      font-size: 0.85rem;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .sidebar li a:hover {
      background: var(--pico-muted-border-color);
    }
    .sidebar li a.active {
      background: var(--pico-primary-background);
      color: var(--pico-primary-inverse);
    }
    main {
      flex: 1;
      min-width: 0;
      padding: 2rem 3rem;
      max-width: 960px;
    }
    @media (max-width: 768px) {
      body { flex-direction: column; }
      .sidebar {
        width: 100%;
        height: auto;
        position: static;
        border-right: none;
        border-bottom: 1px solid var(--pico-muted-border-color);
      }
      main { padding: 1rem; }
    }
"""


def render_page(body_html: str, title: str, current_path: str) -> str:
    sidebar = render_sidebar(current_path)
    return f"""<!doctype html>
<html lang="en" data-theme="light">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{title}</title>
  <link rel="stylesheet" href="{PICO_CSS}">
  <style>{STYLE}</style>
</head>
<body>
{sidebar}
<main>
{body_html}
</main>
</body>
</html>"""


def render_index() -> str:
    """Render the home page with a table of contents."""
    folders = collect_md_files(CONTENT_DIR)
    parts: list[str] = ["<h1>Table of Contents</h1>"]
    for folder, files in folders.items():
        heading = folder if folder else "Files"
        parts.append(f"<h2>{heading}</h2>")
        parts.append('<table role="grid"><thead><tr><th>File</th><th>Title</th></tr></thead><tbody>')
        for fname, title, url_path in files:
            parts.append(
                f'<tr><td><a href="{url_path}">{fname}</a></td>'
                f"<td>{title}</td></tr>"
            )
        parts.append("</tbody></table>")
    body = "\n".join(parts)
    return render_page(body, "Table of Contents", "/")


@app.get("/", response_class=HTMLResponse)
async def index():
    return render_index()


@app.get("/{file_path:path}", response_class=HTMLResponse)
async def serve_md(request: Request, file_path: str):
    target = (CONTENT_DIR / file_path).resolve()

    # Path traversal protection
    if not str(target).startswith(str(CONTENT_DIR.resolve())):
        return HTMLResponse(render_page("<h1>403 Forbidden</h1>", "403", request.url.path), status_code=403)

    if not target.is_file() or target.suffix != ".md":
        return HTMLResponse(render_page("<h1>404 - Not Found</h1>", "404", request.url.path), status_code=404)

    text = read_file(target)
    body_html = markdown.markdown(text, extensions=MD_EXTENSIONS)
    title = extract_title(target)
    return render_page(body_html, title, "/" + file_path)


def main() -> None:
    global CONTENT_DIR
    parser = argparse.ArgumentParser(description="Serve .md files as HTML.")
    parser.add_argument(
        "content_dir",
        nargs="?",
        default=None,
        help="Directory to scan for .md files (default: current directory)",
    )
    parser.add_argument(
        "--port",
        type=int,
        default=int(os.environ.get("PORT", "8000")),
        help="Port to listen on (default: 8000 or $PORT)",
    )
    args = parser.parse_args()

    if args.content_dir:
        CONTENT_DIR = Path(args.content_dir).resolve()

    host = "127.0.0.1"
    print(f"Serving .md files from: {CONTENT_DIR}")
    print(f"Open http://{host}:{args.port}")
    uvicorn.run(app, host=host, port=args.port, log_level="info")


if __name__ == "__main__":
    main()
