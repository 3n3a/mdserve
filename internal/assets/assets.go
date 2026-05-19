package assets

import "embed"

//go:embed pico.min.css mermaid.min.js
var FS embed.FS

func Empty() bool {
	for _, name := range []string{"pico.min.css", "mermaid.min.js"} {
		data, err := FS.ReadFile(name)
		if err != nil || len(data) == 0 {
			return true
		}
	}
	return false
}
