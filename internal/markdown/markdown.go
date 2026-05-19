package markdown

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/3n3a/mdserve/internal/collector"
	"github.com/3n3a/mdserve/internal/fileio"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

var (
	listRE         = regexp.MustCompile(`^(\s*)([-*+]|\d+\.)\s`)
	mermaidBlockRE = regexp.MustCompile("(?ms)^```mermaid[ \\t]*\\r?\\n(.*?)\\r?\\n```[ \\t]*$")

	md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
			highlighting.NewHighlighting(
				highlighting.WithStyle("github"),
			),
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

	text, mermaids := extractMermaid(text)

	var buf bytes.Buffer
	if err := md.Convert([]byte(text), &buf); err != nil {
		return "", "", err
	}
	body := reinjectMermaid(buf.String(), mermaids)
	return template.HTML(body), collector.ExtractTitle(path), nil
}

func extractMermaid(text string) (string, []string) {
	var mermaids []string
	out := mermaidBlockRE.ReplaceAllStringFunc(text, func(m string) string {
		sub := mermaidBlockRE.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		idx := len(mermaids)
		mermaids = append(mermaids, sub[1])
		return fmt.Sprintf("\n\n@@MDSERVE_MERMAID_%d@@\n\n", idx)
	})
	return out, mermaids
}

func reinjectMermaid(body string, mermaids []string) string {
	for i, diag := range mermaids {
		placeholder := fmt.Sprintf("<p>@@MDSERVE_MERMAID_%d@@</p>", i)
		replacement := `<pre class="mermaid">` + diag + `</pre>`
		body = strings.Replace(body, placeholder, replacement, 1)
	}
	return body
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
