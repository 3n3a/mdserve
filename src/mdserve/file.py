from pathlib import Path

def read_file(path: Path) -> str:
    """Read a file trying multiple encodings."""
    for encoding in ("utf-8-sig", "utf-8", "latin-1"):
        try:
            return path.read_text(encoding=encoding)
        except (UnicodeDecodeError, ValueError):
            continue
    return path.read_text(encoding="latin-1", errors="replace")
