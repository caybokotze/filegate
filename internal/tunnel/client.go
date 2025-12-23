package tunnel

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/filegate/filegate/internal/protocol"
	"github.com/gorilla/websocket"
)

const (
	// Version is the client protocol version
	Version = "1.0.0"

	// Reconnect settings
	initialReconnectDelay = 1 * time.Second
	maxReconnectDelay     = 30 * time.Second
	reconnectMultiplier   = 2.0

	// Ping/pong settings
	pingInterval = 30 * time.Second
	pongWait     = 10 * time.Second
)

// Client manages the WebSocket connection to the relay server
type Client struct {
	relayURL string
	handler  http.Handler
	conn     *websocket.Conn
	mu       sync.Mutex

	subdomain string
	fullURL   string

	onConnected    func(subdomain, fullURL string)
	onDisconnected func(err error)
	onReconnecting func(attempt int)
}

// Config holds configuration for the tunnel client
type Config struct {
	// RelayURL is the WebSocket URL of the relay server (e.g., "wss://davproxy.com/tunnel")
	RelayURL string
	// Handler is the HTTP handler (WebDAV server) to forward requests to
	Handler http.Handler
	// OnConnected is called when connection is established
	OnConnected func(subdomain, fullURL string)
	// OnDisconnected is called when connection is lost
	OnDisconnected func(err error)
	// OnReconnecting is called when attempting to reconnect
	OnReconnecting func(attempt int)
}

// New creates a new tunnel client
func New(cfg Config) *Client {
	return &Client{
		relayURL:       cfg.RelayURL,
		handler:        cfg.Handler,
		onConnected:    cfg.OnConnected,
		onDisconnected: cfg.OnDisconnected,
		onReconnecting: cfg.OnReconnecting,
	}
}

// Connect establishes a connection to the relay server and blocks until closed
func (c *Client) Connect(ctx context.Context) error {
	return c.connectWithRetry(ctx)
}

func (c *Client) connectWithRetry(ctx context.Context) error {
	delay := initialReconnectDelay
	attempt := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := c.connectOnce(ctx)
		if err == nil {
			// Connection closed cleanly
			return nil
		}

		if c.onDisconnected != nil {
			c.onDisconnected(err)
		}

		attempt++
		if c.onReconnecting != nil {
			c.onReconnecting(attempt)
		}

		// Wait before reconnecting
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		// Increase delay for next attempt (exponential backoff)
		delay = time.Duration(float64(delay) * reconnectMultiplier)
		if delay > maxReconnectDelay {
			delay = maxReconnectDelay
		}
	}
}

func (c *Client) connectOnce(ctx context.Context) error {
	// Connect to relay
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, c.relayURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()

	defer func() {
		conn.Close()
		c.mu.Lock()
		c.conn = nil
		c.mu.Unlock()
	}()

	// Send registration
	if err := c.register(); err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	// Wait for registration confirmation
	if err := c.waitForRegistered(); err != nil {
		return fmt.Errorf("registration confirmation failed: %w", err)
	}

	// Notify connected
	if c.onConnected != nil {
		c.onConnected(c.subdomain, c.fullURL)
	}

	// Start ping goroutine
	go c.pingLoop(ctx)

	// Handle messages
	return c.handleMessages(ctx)
}

func (c *Client) register() error {
	msg, err := protocol.NewMessage(protocol.TypeRegister, protocol.RegisterPayload{
		Version: Version,
	})
	if err != nil {
		return err
	}

	data, err := msg.Marshal()
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *Client) waitForRegistered() error {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()

	// Set read deadline for registration response
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer conn.SetReadDeadline(time.Time{})

	_, data, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	msg, err := protocol.Unmarshal(data)
	if err != nil {
		return err
	}

	if msg.Type == protocol.TypeError {
		var errPayload protocol.ErrorPayload
		msg.ParsePayload(&errPayload)
		return fmt.Errorf("server error: %s", errPayload.Message)
	}

	if msg.Type != protocol.TypeRegistered {
		return fmt.Errorf("unexpected message type: %s", msg.Type)
	}

	var payload protocol.RegisteredPayload
	if err := msg.ParsePayload(&payload); err != nil {
		return err
	}

	c.subdomain = payload.Subdomain
	c.fullURL = payload.FullURL

	return nil
}

func (c *Client) handleMessages(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		c.mu.Lock()
		conn := c.conn
		c.mu.Unlock()

		if conn == nil {
			return fmt.Errorf("connection closed")
		}

		_, data, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		msg, err := protocol.Unmarshal(data)
		if err != nil {
			continue // Skip malformed messages
		}

		switch msg.Type {
		case protocol.TypeHTTPRequest:
			go c.handleHTTPRequest(msg)
		case protocol.TypePong:
			// Pong received, connection is healthy
		case protocol.TypeError:
			var errPayload protocol.ErrorPayload
			msg.ParsePayload(&errPayload)
			return fmt.Errorf("server error: %s", errPayload.Message)
		}
	}
}

func (c *Client) handleHTTPRequest(msg *protocol.Message) {
	var reqPayload protocol.HTTPRequestPayload
	if err := msg.ParsePayload(&reqPayload); err != nil {
		return
	}

	// Create HTTP request
	req := httptest.NewRequest(reqPayload.Method, reqPayload.Path, bytes.NewReader(reqPayload.Body))
	for key, values := range reqPayload.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Create response recorder
	rec := httptest.NewRecorder()

	// Handle the request with WebDAV server
	c.handler.ServeHTTP(rec, req)

	// Read response body
	respBody, _ := io.ReadAll(rec.Body)

	// Send response back
	respPayload := protocol.HTTPResponsePayload{
		ID:         reqPayload.ID,
		StatusCode: rec.Code,
		Headers:    rec.Header(),
		Body:       respBody,
	}

	respMsg, err := protocol.NewMessage(protocol.TypeHTTPResponse, respPayload)
	if err != nil {
		return
	}

	data, err := respMsg.Marshal()
	if err != nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (c *Client) pingLoop(ctx context.Context) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.sendPing()
		}
	}
}

func (c *Client) sendPing() {
	msg, err := protocol.NewMessage(protocol.TypePing, nil)
	if err != nil {
		return
	}

	data, err := msg.Marshal()
	if err != nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		c.conn.WriteMessage(websocket.TextMessage, data)
	}
}

// Close closes the connection
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Subdomain returns the assigned subdomain
func (c *Client) Subdomain() string {
	return c.subdomain
}

// FullURL returns the full URL for accessing the WebDAV
func (c *Client) FullURL() string {
	return c.fullURL
}
