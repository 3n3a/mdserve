from __future__ import annotations

import argparse

import os
from pathlib import Path

import uvicorn
from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse

from mdserve.render import *
from mdserve.markdown import *

CONTENT_DIR: Path = Path.cwd()

app = FastAPI()

@app.get("/", response_class=HTMLResponse)
async def index():
    return render_index(CONTENT_DIR)


@app.get("/{file_path:path}", response_class=HTMLResponse)
async def serve_md(request: Request, file_path: str):
    target = (CONTENT_DIR / file_path).resolve()

    # Path traversal protection
    if not str(target).startswith(str(CONTENT_DIR.resolve())):
        return HTMLResponse(render_page(CONTENT_DIR, "<h1>403 Forbidden</h1>", "403", request.url.path), status_code=403)

    if not target.is_file() or target.suffix != ".md":
        return HTMLResponse(render_page(CONTENT_DIR, "<h1>404 - Not Found</h1>", "404", request.url.path), status_code=404)

    body_html, title = convert_to_markdown(target)
    return render_page(CONTENT_DIR, body_html, title, "/" + file_path)


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
