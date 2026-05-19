package render

import (
	"html/template"
	"io"

	"github.com/3n3a/mdserve/internal/collector"
)

type AssetURLs struct {
	Pico    string
	Mermaid string
}

type PageData struct {
	Title       string
	Body        template.HTML
	CurrentPath string
	Groups      []collector.Group
	Style       template.CSS
	PicoHref    string
	MermaidHref string
}

type indexData struct {
	Groups []collector.Group
}

var (
	pageT  = template.Must(template.New("page").Parse(pageTmpl))
	indexT = template.Must(template.New("index").Parse(indexTmpl))
)

func RenderPage(w io.Writer, title string, body template.HTML, currentPath string, groups []collector.Group, urls AssetURLs) error {
	return pageT.Execute(w, PageData{
		Title:       title,
		Body:        body,
		CurrentPath: currentPath,
		Groups:      groups,
		Style:       template.CSS(style),
		PicoHref:    urls.Pico,
		MermaidHref: urls.Mermaid,
	})
}

func RenderIndex(w io.Writer, groups []collector.Group, urls AssetURLs) error {
	var body templateBuffer
	if err := indexT.Execute(&body, indexData{Groups: groups}); err != nil {
		return err
	}
	return RenderPage(w, "Table of Contents", template.HTML(body.String()), "/", groups, urls)
}

type templateBuffer struct {
	b []byte
}

func (t *templateBuffer) Write(p []byte) (int, error) { t.b = append(t.b, p...); return len(p), nil }
func (t *templateBuffer) String() string              { return string(t.b) }
