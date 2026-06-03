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
	ContentDir     string
	User           string
	Pass           string
	Offline        bool
	SaveCheckboxes bool
}

type Server struct {
	opts        Options
	contentRoot string
	secret      []byte
}

func New(opts Options) (http.Handler, error) {
	abs, err := filepath.Abs(opts.ContentDir)
	if err != nil {
		return nil, err
	}
	secret, err := newSessionSecret()
	if err != nil {
		return nil, err
	}
	s := &Server{opts: opts, contentRoot: abs, secret: secret}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/login", s.handleLoginGet)
	r.Post("/login", s.handleLoginPost)
	r.Get("/logout", s.handleLogout)
	r.Handle("/_assets/*", http.StripPrefix("/_assets/", http.FileServer(http.FS(assets.FS))))

	r.Group(func(r chi.Router) {
		if opts.User != "" {
			r.Use(s.requireAuth)
		}
		r.Get("/", s.handleIndex)
		r.Get("/*", s.handleFile)
	})
	return r, nil
}

func (s *Server) assetURLs() render.AssetURLs {
	if s.opts.Offline {
		return render.AssetURLs{Pico: picoLocal, Mermaid: mermaidLcl}
	}
	return render.AssetURLs{Pico: picoCDN, Mermaid: mermaidCDN}
}
