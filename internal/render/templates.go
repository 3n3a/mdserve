package render

const style = `
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
`

const pageTmpl = `<!doctype html>
<html lang="en" data-theme="light">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}}</title>
  <link rel="stylesheet" href="{{.PicoHref}}">
  <style>{{.Style}}</style>
  <script src="{{.MermaidHref}}"></script>
  <script>mermaid.initialize({ startOnLoad: true });</script>
</head>
<body>
<aside class="sidebar">
<a href="/" class="sidebar-logo"><strong>mdserve</strong></a>
{{- range .Groups}}
  <div class="sidebar-heading">{{if .Folder}}{{.Folder}}{{else}}Files{{end}}</div>
  <ul>
  {{- range .Entries}}
    <li><a href="{{.URLPath}}"{{if eq .URLPath $.CurrentPath}} class="active"{{end}}>{{.Title}}</a></li>
  {{- end}}
  </ul>
{{- end}}
</aside>
<main>
{{.Body}}
</main>
</body>
</html>`

const indexTmpl = `<h1>Table of Contents</h1>
{{- range .Groups}}
<h2>{{if .Folder}}{{.Folder}}{{else}}Files{{end}}</h2>
<table role="grid"><thead><tr><th>File</th><th>Title</th></tr></thead><tbody>
{{- range .Entries}}
<tr><td><a href="{{.URLPath}}">{{.File}}</a></td><td>{{.Title}}</td></tr>
{{- end}}
</tbody></table>
{{- end}}`
