package server

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/3n3a/mdserve/internal/render"
)

const (
	sessionCookieName = "mdserve_session"
	sessionTTL        = 24 * time.Hour
)

func newSessionSecret() ([]byte, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, err
	}
	return secret, nil
}

func (s *Server) issueSession() string {
	expiry := time.Now().Add(sessionTTL).Unix()
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(expiry))
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(buf)
	return base64.RawURLEncoding.EncodeToString(append(buf, mac.Sum(nil)...))
}

func (s *Server) validSession(token string) bool {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(raw) != 8+sha256.Size {
		return false
	}
	expiryBytes := raw[:8]
	got := raw[8:]
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(expiryBytes)
	if !hmac.Equal(got, mac.Sum(nil)) {
		return false
	}
	expiry := int64(binary.BigEndian.Uint64(expiryBytes))
	return time.Now().Unix() < expiry
}

func (s *Server) isAuthenticated(r *http.Request) bool {
	if s.opts.User == "" {
		return true
	}
	c, err := r.Cookie(sessionCookieName)
	if err != nil {
		return false
	}
	return s.validSession(c.Value)
}

func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.isAuthenticated(r) {
			next.ServeHTTP(w, r)
			return
		}
		nextURL := r.URL.RequestURI()
		http.Redirect(w, r, "/login?next="+url.QueryEscape(nextURL), http.StatusSeeOther)
	})
}

func (s *Server) handleLoginGet(w http.ResponseWriter, r *http.Request) {
	if s.opts.User == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if s.isAuthenticated(r) {
		http.Redirect(w, r, safeNext(r.URL.Query().Get("next")), http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = render.RenderLogin(w, "", r.URL.Query().Get("next"), s.assetURLs())
}

func (s *Server) handleLoginPost(w http.ResponseWriter, r *http.Request) {
	if s.opts.User == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if err := r.ParseForm(); err != nil {
		s.renderLoginError(w, "Bad request", "")
		return
	}
	user := strings.TrimSpace(r.FormValue("user"))
	pass := r.FormValue("pass")
	next := r.FormValue("next")

	expectedUser := strings.TrimSpace(s.opts.User)
	expectedPass := strings.TrimSpace(s.opts.Pass)

	if subtle.ConstantTimeCompare([]byte(user), []byte(expectedUser)) != 1 ||
		subtle.ConstantTimeCompare([]byte(pass), []byte(expectedPass)) != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		s.renderLoginError(w, "Invalid credentials", next)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    s.issueSession(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(sessionTTL.Seconds()),
	})
	http.Redirect(w, r, safeNext(next), http.StatusSeeOther)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (s *Server) renderLoginError(w http.ResponseWriter, msg, next string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = render.RenderLogin(w, msg, next, s.assetURLs())
}

func safeNext(next string) string {
	if next == "" || !strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//") {
		return "/"
	}
	return next
}
