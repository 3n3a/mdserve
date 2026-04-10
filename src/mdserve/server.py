"""Minimal webserver that converts .md files to HTML on the fly using Pico CSS."""

from __future__ import annotations

import argparse
import os
import re
import urllib.parse
from http.server import HTTPServer, BaseHTTPRequestHandler
from pathlib import Path

import markdown

# The directory containing .md files (one level up from the repo root).
# Layout: src/mdserve/server.py  →  .parent = src/mdserve  →  .parent = src
#         →  .parent = repo root  →  .parent = content dir
_DEFAULT_CONTENT_DIR = Path(__file__).resolve().parent.parent.parent.parent
CONTENT_DIR: Path = _DEFAULT_CONTENT_DIR

PICO_CSS = "https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css"


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


def render_sidebar(current_path: str) -> str:
    """Build the sidebar navigation HTML."""
    folders = collect_md_files(CONTENT_DIR)
    parts: list[str] = []
    parts.append('<nav id="sidebar">')
    parts.append('<a href="/" class="logo"><strong>MD Viewer</strong></a>')
    for folder, files in folders.items():
        heading = folder if folder else "Root"
        parts.append(f"<details open><summary>{heading}</summary><ul>")
        for _fname, title, url_path in files:
            active = ' class="active"' if url_path == current_path else ""
            parts.append(
                f'<li><a href="{url_path}"{active}>{title}</a></li>'
            )
        parts.append("</ul></details>")
    parts.append("</nav>")
    return "\n".join(parts)


def render_page(body_html: str, title: str, current_path: str) -> str:
    sidebar = render_sidebar(current_path)
    return f"""<!doctype html>
<html lang="de" data-theme="light">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{title}</title>
  <link rel="stylesheet" href="{PICO_CSS}">
  <style>
    body {{
      display: flex;
      min-height: 100vh;
      margin: 0;
    }}
    #sidebar {{
      width: 320px;
      min-width: 260px;
      padding: 1.5rem 1rem;
      border-right: 1px solid var(--pico-muted-border-color);
      overflow-y: auto;
      position: sticky;
      top: 0;
      height: 100vh;
      font-size: 0.9rem;
    }}
    #sidebar .logo {{
      display: block;
      margin-bottom: 1rem;
      font-size: 1.1rem;
      text-decoration: none;
    }}
    #sidebar details {{
      margin-bottom: 0.25rem;
      border: none;
    }}
    #sidebar summary {{
      font-weight: 600;
      padding: 0.3rem 0;
      font-size: 0.85rem;
      text-transform: uppercase;
      letter-spacing: 0.04em;
      color: var(--pico-muted-color);
    }}
    #sidebar ul {{
      list-style: none;
      padding-left: 0.5rem;
      margin: 0;
    }}
    #sidebar li {{
      margin: 0;
      padding: 0;
    }}
    #sidebar a {{
      display: block;
      padding: 0.25rem 0.5rem;
      border-radius: 4px;
      text-decoration: none;
      color: var(--pico-color);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }}
    #sidebar a:hover,
    #sidebar a.active {{
      background: var(--pico-primary-focus);
    }}
    main {{
      flex: 1;
      padding: 2rem 3rem;
      max-width: 960px;
      overflow-x: auto;
    }}
    main table {{
      font-size: 0.88rem;
    }}
    @media (max-width: 768px) {{
      body {{ flex-direction: column; }}
      #sidebar {{
        width: 100%;
        height: auto;
        position: static;
        border-right: none;
        border-bottom: 1px solid var(--pico-muted-border-color);
      }}
      main {{ padding: 1rem; }}
    }}
  </style>
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
    parts: list[str] = ["<h1>Inhaltsverzeichnis</h1>"]
    for folder, files in folders.items():
        heading = folder if folder else "Stamm-Verzeichnis"
        parts.append(f"<h2>{heading}</h2>")
        parts.append('<table role="grid"><thead><tr><th>Datei</th><th>Titel</th></tr></thead><tbody>')
        for fname, title, url_path in files:
            parts.append(
                f'<tr><td><a href="{url_path}">{fname}</a></td>'
                f"<td>{title}</td></tr>"
            )
        parts.append("</tbody></table>")
    body = "\n".join(parts)
    return render_page(body, "Inhaltsverzeichnis", "/")


MD_EXTENSIONS = ["tables", "fenced_code", "toc", "sane_lists", "smarty"]


class Handler(BaseHTTPRequestHandler):
    def do_GET(self) -> None:
        parsed = urllib.parse.urlparse(self.path)
        req_path = urllib.parse.unquote(parsed.path)

        if req_path == "/" or req_path == "":
            html = render_index()
            self._respond(200, html)
            return

        # Map URL path to file
        rel = req_path.lstrip("/")
        file_path = CONTENT_DIR / rel

        if file_path.is_file() and file_path.suffix == ".md":
            try:
                text = file_path.read_text(encoding="utf-8")
            except Exception as exc:
                self._respond(500, render_page(f"<p>Error reading file: {exc}</p>", "Error", req_path))
                return
            body_html = markdown.markdown(text, extensions=MD_EXTENSIONS)
            title = extract_title(file_path)
            html = render_page(body_html, title, req_path)
            self._respond(200, html)
        else:
            html = render_page("<h1>404 - Nicht gefunden</h1>", "404", req_path)
            self._respond(404, html)

    def _respond(self, code: int, html: str) -> None:
        payload = html.encode("utf-8")
        self.send_response(code)
        self.send_header("Content-Type", "text/html; charset=utf-8")
        self.send_header("Content-Length", str(len(payload)))
        self.end_headers()
        self.wfile.write(payload)

    def log_message(self, fmt: str, *args) -> None:
        print(f"[{self.log_date_time_string()}] {fmt % args}")


def main() -> None:
    global CONTENT_DIR

    parser = argparse.ArgumentParser(description="Serve .md files as HTML.")
    parser.add_argument(
        "content_dir",
        nargs="?",
        type=Path,
        default=_DEFAULT_CONTENT_DIR,
        help="Directory containing .md files (default: parent of this script)",
    )
    parser.add_argument(
        "--port",
        type=int,
        default=int(os.environ.get("PORT", "8000")),
        help="Port to listen on (default: 8000 or $PORT)",
    )
    args = parser.parse_args()

    CONTENT_DIR = args.content_dir.resolve()

    host = "127.0.0.1"
    server = HTTPServer((host, args.port), Handler)
    print(f"Serving .md files from: {CONTENT_DIR}")
    print(f"Open http://{host}:{args.port}")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\nShutting down.")
        server.server_close()


if __name__ == "__main__":
    main()
