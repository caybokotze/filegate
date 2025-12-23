package relay

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/filegate/filegate/internal/protocol"
	"github.com/gorilla/websocket"
)

const (
	// MaxSubdomainAttempts is the maximum number of attempts to generate a unique subdomain
	MaxSubdomainAttempts = 10
	// RequestTimeout is the maximum time to wait for a response from the client
	RequestTimeout = 60 * time.Second
	// WriteTimeout is the timeout for writing messages to the client
	WriteTimeout = 10 * time.Second
)

// Client represents a connected tunnel client
type Client struct {
	subdomain string
	conn      *websocket.Conn
	mu        sync.Mutex

	// pending tracks pending requests waiting for responses
	pending   map[string]chan *protocol.HTTPResponsePayload
	pendingMu sync.Mutex
}

// Hub manages all connected tunnel clients
type Hub struct {
	clients map[string]*Client
	mu      sync.RWMutex

	domain string // Base domain (e.g., "davproxy.com")
}

// NewHub creates a new hub
func NewHub(domain string) *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		domain:  domain,
	}
}

// Register adds a new client and returns the assigned subdomain
func (h *Hub) Register(conn *websocket.Conn) (*Client, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Generate unique subdomain
	var subdomain string
	for i := 0; i < MaxSubdomainAttempts; i++ {
		var err error
		subdomain, err = GenerateSubdomain()
		if err != nil {
			return nil, fmt.Errorf("failed to generate subdomain: %w", err)
		}

		if _, exists := h.clients[subdomain]; !exists {
			break
		}

		if i == MaxSubdomainAttempts-1 {
			return nil, fmt.Errorf("failed to generate unique subdomain after %d attempts", MaxSubdomainAttempts)
		}
	}

	client := &Client{
		subdomain: subdomain,
		conn:      conn,
		pending:   make(map[string]chan *protocol.HTTPResponsePayload),
	}

	h.clients[subdomain] = client
	return client, nil
}

// Unregister removes a client
func (h *Hub) Unregister(subdomain string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client, exists := h.clients[subdomain]; exists {
		// Cancel all pending requests
		client.pendingMu.Lock()
		for _, ch := range client.pending {
			close(ch)
		}
		client.pending = make(map[string]chan *protocol.HTTPResponsePayload)
		client.pendingMu.Unlock()

		delete(h.clients, subdomain)
	}
}

// GetClient returns a client by subdomain
func (h *Hub) GetClient(subdomain string) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients[subdomain]
}

// Domain returns the base domain
func (h *Hub) Domain() string {
	return h.domain
}

// ClientCount returns the number of connected clients
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// Subdomain returns the client's subdomain
func (c *Client) Subdomain() string {
	return c.subdomain
}

// SendRequest sends an HTTP request to the client and waits for a response
func (c *Client) SendRequest(ctx context.Context, req *protocol.HTTPRequestPayload) (*protocol.HTTPResponsePayload, error) {
	// Create response channel
	respChan := make(chan *protocol.HTTPResponsePayload, 1)

	c.pendingMu.Lock()
	c.pending[req.ID] = respChan
	c.pendingMu.Unlock()

	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, req.ID)
		c.pendingMu.Unlock()
	}()

	// Send request
	msg, err := protocol.NewMessage(protocol.TypeHTTPRequest, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	data, err := msg.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	c.mu.Lock()
	c.conn.SetWriteDeadline(time.Now().Add(WriteTimeout))
	err = c.conn.WriteMessage(websocket.TextMessage, data)
	c.mu.Unlock()

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Wait for response
	select {
	case resp, ok := <-respChan:
		if !ok {
			return nil, fmt.Errorf("connection closed")
		}
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(RequestTimeout):
		return nil, fmt.Errorf("request timeout")
	}
}

// HandleResponse processes a response from the client
func (c *Client) HandleResponse(resp *protocol.HTTPResponsePayload) {
	c.pendingMu.Lock()
	defer c.pendingMu.Unlock()

	if ch, exists := c.pending[resp.ID]; exists {
		select {
		case ch <- resp:
		default:
		}
	}
}

// SendPong sends a pong response
func (c *Client) SendPong() error {
	msg, err := protocol.NewMessage(protocol.TypePong, nil)
	if err != nil {
		return err
	}

	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn.SetWriteDeadline(time.Now().Add(WriteTimeout))
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// Close closes the client connection
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.Close()
}
