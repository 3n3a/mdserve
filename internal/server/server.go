package server

import (
	"net/http"
	"path/filepath"

	"github.com/3n3a/mdserve/internal/assets"
	"github.com/3n3a/mdserve/internal/render"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	picoCDN    = "https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css"
	mermaidCDN = "https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.min.js"
	picoLocal  = "/_assets/pico.min.css"
	mermaidLcl = "/_assets/mermaid.min.js"
)

type Options struct {
	ContentDir string
	User       string
	Pass       string
	Offline    bool
}

type Server struct {
	opts        Options
	contentRoot string
}

func New(opts Options) (http.Handler, error) {
	abs, err := filepath.Abs(opts.ContentDir)
	if err != nil {
		return nil, err
	}
	s := &Server{opts: opts, contentRoot: abs}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	if opts.User != "" {
		r.Use(basicAuth(opts.User, opts.Pass))
	}

	r.Get("/", s.handleIndex)
	r.Handle("/_assets/*", http.StripPrefix("/_assets/", http.FileServer(http.FS(assets.FS))))
	r.Get("/*", s.handleFile)
	return r, nil
}

func (s *Server) assetURLs() render.AssetURLs {
	if s.opts.Offline {
		return render.AssetURLs{Pico: picoLocal, Mermaid: mermaidLcl}
	}
	return render.AssetURLs{Pico: picoCDN, Mermaid: mermaidCDN}
}
