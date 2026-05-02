from pathlib import Path
from mdserve.collector import *

PICO_CSS = "https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css"
MERMAID_JS = "https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs"

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

def render_index(content_dir: Path) -> str:
    """Render the home page with a table of contents."""
    folders = collect_md_files(content_dir)
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
    return render_page(content_dir, body, "Table of Contents", "/")

def render_page(content_dir: Path, body_html: str, title: str, current_path: str) -> str:
    sidebar = render_sidebar(content_dir, current_path)
    return f"""<!doctype html>
<html lang="en" data-theme="light">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{title}</title>
  <link rel="stylesheet" href="{PICO_CSS}">
  <style>{STYLE}</style>
  <script type="module">
    import mermaid from '{MERMAID_JS}';
    mermaid.initialize({{ startOnLoad: true }});
  </script>
</head>
<body>
{sidebar}
<main>
{body_html}
</main>
</body>
</html>"""


def render_sidebar(content_dir: Path, current_path: str) -> str:
    """Build the sidebar navigation HTML."""
    folders = collect_md_files(content_dir)
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