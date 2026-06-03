package server

import (
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/3n3a/mdserve/internal/collector"
	"github.com/3n3a/mdserve/internal/markdown"
	"github.com/3n3a/mdserve/internal/render"
	"github.com/go-chi/chi/v5"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	groups, err := collector.Collect(s.opts.ContentDir)
	if err != nil {
		s.errorPage(w, r, http.StatusInternalServerError, "500 - "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := render.RenderIndex(w, groups, s.assetURLs()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleFile(w http.ResponseWriter, r *http.Request) {
	reqPath := chi.URLParam(r, "*")
	cleaned := filepath.Clean("/" + reqPath)
	target := filepath.Join(s.contentRoot, filepath.FromSlash(cleaned))

	rel, err := filepath.Rel(s.contentRoot, target)
	if err != nil || strings.HasPrefix(rel, "..") || rel == ".." {
		s.errorPage(w, r, http.StatusForbidden, "403 Forbidden")
		return
	}

	info, err := os.Stat(target)
	if err != nil || info.IsDir() || !strings.EqualFold(filepath.Ext(target), ".md") {
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			s.errorPage(w, r, http.StatusInternalServerError, "500 - "+err.Error())
			return
		}
		s.errorPage(w, r, http.StatusNotFound, "404 - Not Found")
		return
	}

	body, title, err := markdown.Convert(target)
	if err != nil {
		s.errorPage(w, r, http.StatusInternalServerError, "500 - "+err.Error())
		return
	}

	groups, err := collector.Collect(s.opts.ContentDir)
	if err != nil {
		s.errorPage(w, r, http.StatusInternalServerError, "500 - "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := render.RenderPage(w, title, body, "/"+filepath.ToSlash(rel), groups, s.assetURLs(), s.opts.SaveCheckboxes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) errorPage(w http.ResponseWriter, _ *http.Request, status int, message string) {
	groups, _ := collector.Collect(s.opts.ContentDir)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_ = render.RenderPage(w, http.StatusText(status), template.HTML("<h1>"+template.HTMLEscapeString(message)+"</h1>"), "", groups, s.assetURLs(), false)
}
