package webdav

import (
	"crypto/subtle"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
)

// Server wraps a WebDAV handler with authentication
type Server struct {
	handler  *webdav.Handler
	username string
	password string
}

// Config holds configuration for the WebDAV server
type Config struct {
	// Root directory to serve (defaults to current working directory)
	Root string
	// Username for Basic Auth
	Username string
	// Password for Basic Auth
	Password string
}

// New creates a new WebDAV server
func New(cfg Config) (*Server, error) {
	root := cfg.Root
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	handler := &webdav.Handler{
		FileSystem: webdav.Dir(root),
		LockSystem: webdav.NewMemLS(),
		Prefix:     "",
	}

	return &Server{
		handler:  handler,
		username: cfg.Username,
		password: cfg.Password,
	}, nil
}

// ServeHTTP implements http.Handler with Basic Auth
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check Basic Auth
	if !s.authenticate(r) {
		w.Header().Set("WWW-Authenticate", `Basic realm="filegate"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	s.handler.ServeHTTP(w, r)
}

// authenticate checks the request for valid Basic Auth credentials
func (s *Server) authenticate(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false
	}

	// Use constant-time comparison to prevent timing attacks
	usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(s.username)) == 1
	passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(s.password)) == 1

	return usernameMatch && passwordMatch
}

// Handler returns the underlying http.Handler for use with custom servers
func (s *Server) Handler() http.Handler {
	return s
}
