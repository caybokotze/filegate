package relay

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/filegate/filegate/internal/protocol"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for WebSocket
	},
}

// Server is the relay server that handles tunnel connections and HTTP proxying
type Server struct {
	hub    *Hub
	mux    *http.ServeMux
	domain string
}

// Config holds configuration for the relay server
type Config struct {
	// Domain is the base domain (e.g., "davproxy.com")
	Domain string
	// Port is the port to listen on
	Port int
}

// NewServer creates a new relay server
func NewServer(cfg Config) *Server {
	hub := NewHub(cfg.Domain)
	s := &Server{
		hub:    hub,
		mux:    http.NewServeMux(),
		domain: cfg.Domain,
	}

	// Register routes
	s.mux.HandleFunc("/tunnel", s.handleTunnel)
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/", s.handleProxy)

	return s
}

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// handleHealth returns server health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"ok","clients":%d}`, s.hub.ClientCount())
}

// handleTunnel handles WebSocket connections from CLI clients
func (s *Server) handleTunnel(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// Wait for registration message
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, data, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Failed to read registration: %v", err)
		conn.Close()
		return
	}
	conn.SetReadDeadline(time.Time{})

	msg, err := protocol.Unmarshal(data)
	if err != nil || msg.Type != protocol.TypeRegister {
		s.sendError(conn, "invalid_registration", "Expected register message")
		conn.Close()
		return
	}

	// Register client
	client, err := s.hub.Register(conn)
	if err != nil {
		s.sendError(conn, "registration_failed", err.Error())
		conn.Close()
		return
	}

	log.Printf("Client registered: %s", client.Subdomain())

	// Send registration confirmation
	fullURL := fmt.Sprintf("https://%s.%s", client.Subdomain(), s.domain)
	regPayload := protocol.RegisteredPayload{
		Subdomain: client.Subdomain(),
		FullURL:   fullURL,
	}

	respMsg, _ := protocol.NewMessage(protocol.TypeRegistered, regPayload)
	respData, _ := respMsg.Marshal()
	conn.WriteMessage(websocket.TextMessage, respData)

	// Handle messages from client
	defer func() {
		s.hub.Unregister(client.Subdomain())
		conn.Close()
		log.Printf("Client disconnected: %s", client.Subdomain())
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}

		msg, err := protocol.Unmarshal(data)
		if err != nil {
			continue
		}

		switch msg.Type {
		case protocol.TypeHTTPResponse:
			var resp protocol.HTTPResponsePayload
			if err := msg.ParsePayload(&resp); err == nil {
				client.HandleResponse(&resp)
			}
		case protocol.TypePing:
			client.SendPong()
		}
	}
}

// handleProxy handles HTTP requests and proxies them to the appropriate client
func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) {
	// Extract subdomain from Host header
	host := r.Host
	subdomain := s.extractSubdomain(host)

	if subdomain == "" {
		// Request to main domain - show landing page
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><title>filegate</title></head>
<body>
<h1>filegate</h1>
<p>File sharing gateway. Run <code>filegate</code> CLI to expose your files.</p>
</body>
</html>`)
		return
	}

	// Find client
	client := s.hub.GetClient(subdomain)
	if client == nil {
		http.Error(w, "Tunnel not found", http.StatusNotFound)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Create request payload
	reqPayload := &protocol.HTTPRequestPayload{
		ID:      uuid.New().String(),
		Method:  r.Method,
		Path:    r.URL.RequestURI(),
		Headers: r.Header,
		Body:    body,
	}

	// Send request to client and wait for response
	ctx, cancel := context.WithTimeout(r.Context(), RequestTimeout)
	defer cancel()

	resp, err := client.SendRequest(ctx, reqPayload)
	if err != nil {
		log.Printf("Request to %s failed: %v", subdomain, err)
		http.Error(w, "Tunnel error: "+err.Error(), http.StatusBadGateway)
		return
	}

	// Write response
	for key, values := range resp.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(resp.Body)
}

// extractSubdomain extracts the subdomain from a host like "brave-tiger.davproxy.com"
func (s *Server) extractSubdomain(host string) string {
	// Remove port if present
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// Check if it ends with our domain
	suffix := "." + s.domain
	if !strings.HasSuffix(host, suffix) {
		// Check if it's the exact domain
		if host == s.domain {
			return ""
		}
		// Might be localhost or direct IP access during development
		return ""
	}

	// Extract subdomain
	subdomain := strings.TrimSuffix(host, suffix)
	if subdomain == "" || subdomain == host {
		return ""
	}

	return subdomain
}

func (s *Server) sendError(conn *websocket.Conn, code, message string) {
	msg, _ := protocol.NewMessage(protocol.TypeError, protocol.ErrorPayload{
		Code:    code,
		Message: message,
	})
	data, _ := msg.Marshal()
	conn.WriteMessage(websocket.TextMessage, data)
}
