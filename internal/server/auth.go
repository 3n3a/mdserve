package server

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

func (s *Server) basicAuth(next http.Handler) http.Handler {
	expectedUser := []byte(strings.TrimSpace(s.opts.User))
	expectedPass := []byte(strings.TrimSpace(s.opts.Pass))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok ||
			subtle.ConstantTimeCompare([]byte(u), expectedUser) != 1 ||
			subtle.ConstantTimeCompare([]byte(p), expectedPass) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="mdserve"`)
			s.errorPage(w, r, http.StatusUnauthorized, "401 - Authentication required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
