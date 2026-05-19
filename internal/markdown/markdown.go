package markdown

import (
	"bytes"
	"html"
	"html/template"
	"regexp"
	"strings"

	"github.com/3n3a/mdserve/internal/collector"
	"github.com/3n3a/mdserve/internal/fileio"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

var (
	listRE    = regexp.MustCompile(`^(\s*)([-*+]|\d+\.)\s`)
	mermaidRE = regexp.MustCompile(`(?s)<pre><code class="language-mermaid">(.*?)</code></pre>`)

	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
		),
		goldmark.WithRendererOptions(
			gmhtml.WithXHTML(),
		),
	)
)

func Convert(path string) (template.HTML, string, error) {
	text, err := fileio.Read(path)
	if err != nil {
		return "", "", err
	}
	text = ensureBlankLineBeforeLists(text)

	var buf bytes.Buffer
	if err := md.Convert([]byte(text), &buf); err != nil {
		return "", "", err
	}
	body := rewriteMermaid(buf.String())
	return template.HTML(body), collector.ExtractTitle(path), nil
}

func rewriteMermaid(htmlText string) string {
	return mermaidRE.ReplaceAllStringFunc(htmlText, func(m string) string {
		sub := mermaidRE.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		return `<pre class="mermaid">` + html.UnescapeString(sub[1]) + `</pre>`
	})
}

func ensureBlankLineBeforeLists(text string) string {
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))
	for i, line := range lines {
		if i > 0 && listRE.MatchString(line) {
			prev := lines[i-1]
			if strings.TrimSpace(prev) != "" && !listRE.MatchString(prev) {
				out = append(out, "")
			}
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}
