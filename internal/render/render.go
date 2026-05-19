package render

import (
	"bytes"
	"embed"
	"html/template"
	"io"

	"github.com/3n3a/mdserve/internal/collector"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

//go:embed templates/style.css
var styleCSS string

var tmpl = template.Must(template.ParseFS(templatesFS, "templates/*.tmpl"))

type AssetURLs struct {
	Pico    string
	Mermaid string
}

type pageData struct {
	Title       string
	Body        template.HTML
	CurrentPath string
	Groups      []collector.Group
	Style       template.CSS
	PicoHref    string
	MermaidHref string
}

type loginData struct {
	Title    string
	Error    string
	Next     string
	PicoHref string
}

func RenderPage(w io.Writer, title string, body template.HTML, currentPath string, groups []collector.Group, urls AssetURLs) error {
	return tmpl.ExecuteTemplate(w, "page.tmpl", pageData{
		Title:       title,
		Body:        body,
		CurrentPath: currentPath,
		Groups:      groups,
		Style:       template.CSS(styleCSS),
		PicoHref:    urls.Pico,
		MermaidHref: urls.Mermaid,
	})
}

func RenderIndex(w io.Writer, groups []collector.Group, urls AssetURLs) error {
	var body bytes.Buffer
	if err := tmpl.ExecuteTemplate(&body, "index.tmpl", pageData{Groups: groups}); err != nil {
		return err
	}
	return RenderPage(w, "Table of Contents", template.HTML(body.String()), "/", groups, urls)
}

func RenderLogin(w io.Writer, errMsg, next string, urls AssetURLs) error {
	return tmpl.ExecuteTemplate(w, "login.tmpl", loginData{
		Title:    "Sign in — mdserve",
		Error:    errMsg,
		Next:     next,
		PicoHref: urls.Pico,
	})
}
