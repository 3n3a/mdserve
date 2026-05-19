package server

import (
	"crypto/subtle"
	"net/http"
)

func basicAuth(user, pass string) func(http.Handler) http.Handler {
	expectedUser := []byte(user)
	expectedPass := []byte(pass)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if !ok ||
				subtle.ConstantTimeCompare([]byte(u), expectedUser) != 1 ||
				subtle.ConstantTimeCompare([]byte(p), expectedPass) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="mdserve", charset="UTF-8"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
