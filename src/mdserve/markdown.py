import re
import html as _html
import markdown

from mdserve.file import *
from mdserve.collector import *

MD_EXTENSIONS = ["tables", "fenced_code", "toc", "sane_lists", "smarty"]

_LIST_RE = re.compile(r"^(\s*)([-*+]|\d+\.)\s")
_MERMAID_RE = re.compile(
    r'<pre><code class="language-mermaid">(.*?)</code></pre>',
    re.DOTALL,
)

def convert_to_markdown(target: Path):
    text = read_file(target)
    text = _ensure_blank_line_before_lists(text)
    body_html = markdown.markdown(text, extensions=MD_EXTENSIONS)
    body_html = _extract_mermaid_blocks(body_html)
    title = extract_title(target)
    return body_html, title

def _extract_mermaid_blocks(html_text: str) -> str:
    """Replace fenced-code mermaid blocks with <pre class="mermaid"> for Mermaid.js."""
    def _replace(m: re.Match) -> str:
        diagram = _html.unescape(m.group(1))
        return f'<pre class="mermaid">{diagram}</pre>'
    return _MERMAID_RE.sub(_replace, html_text)


def _ensure_blank_line_before_lists(text: str) -> str:
    """Insert a blank line before a list that immediately follows a non-list line.

    Python-markdown requires a blank line between a paragraph and a list;
    without it the list markers are swallowed into the paragraph text.
    """
    lines = text.split("\n")
    result: list[str] = []
    for i, line in enumerate(lines):
        if i > 0 and _LIST_RE.match(line):
            prev = lines[i - 1]
            if prev.strip() and not _LIST_RE.match(prev):
                result.append("")
        result.append(line)
    return "\n".join(result)