import os
from pathlib import Path

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


