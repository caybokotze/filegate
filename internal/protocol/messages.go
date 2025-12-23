package protocol

import (
	"encoding/json"
)

// MessageType defines the type of message being sent
type MessageType string

const (
	// TypeRegister is sent by client to register with the relay
	TypeRegister MessageType = "register"
	// TypeRegistered is sent by relay to confirm registration with assigned subdomain
	TypeRegistered MessageType = "registered"
	// TypeHTTPRequest is sent by relay to forward an HTTP request to the client
	TypeHTTPRequest MessageType = "http_request"
	// TypeHTTPResponse is sent by client to respond to an HTTP request
	TypeHTTPResponse MessageType = "http_response"
	// TypePing is sent to check connection health
	TypePing MessageType = "ping"
	// TypePong is the response to a ping
	TypePong MessageType = "pong"
	// TypeError is sent when an error occurs
	TypeError MessageType = "error"
)

// Message is the base envelope for all WebSocket messages
type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// RegisterPayload is sent by the client to register
type RegisterPayload struct {
	// Version of the client for compatibility checking
	Version string `json:"version"`
}

// RegisteredPayload is sent by the relay after successful registration
type RegisteredPayload struct {
	// Subdomain assigned to this client (e.g., "brave-tiger")
	Subdomain string `json:"subdomain"`
	// FullURL is the complete URL for accessing the WebDAV (e.g., "https://brave-tiger.davproxy.com")
	FullURL string `json:"full_url"`
}

// HTTPRequestPayload represents an incoming HTTP request to be forwarded
type HTTPRequestPayload struct {
	// ID uniquely identifies this request for matching with response
	ID string `json:"id"`
	// Method is the HTTP method (GET, PUT, PROPFIND, etc.)
	Method string `json:"method"`
	// Path is the request path
	Path string `json:"path"`
	// Headers are the HTTP headers
	Headers map[string][]string `json:"headers"`
	// Body is the request body (base64 encoded for binary safety)
	Body []byte `json:"body,omitempty"`
}

// HTTPResponsePayload represents the response to an HTTP request
type HTTPResponsePayload struct {
	// ID matches the request ID
	ID string `json:"id"`
	// StatusCode is the HTTP status code
	StatusCode int `json:"status_code"`
	// Headers are the HTTP response headers
	Headers map[string][]string `json:"headers"`
	// Body is the response body (base64 encoded for binary safety)
	Body []byte `json:"body,omitempty"`
}

// ErrorPayload contains error information
type ErrorPayload struct {
	// Code is a machine-readable error code
	Code string `json:"code"`
	// Message is a human-readable error description
	Message string `json:"message"`
}

// NewMessage creates a new message with the given type and payload
func NewMessage(msgType MessageType, payload interface{}) (*Message, error) {
	var rawPayload json.RawMessage
	if payload != nil {
		var err error
		rawPayload, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	return &Message{
		Type:    msgType,
		Payload: rawPayload,
	}, nil
}

// ParsePayload unmarshals the payload into the provided struct
func (m *Message) ParsePayload(v interface{}) error {
	if m.Payload == nil {
		return nil
	}
	return json.Unmarshal(m.Payload, v)
}

// Marshal converts the message to JSON bytes
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal parses JSON bytes into a Message
func Unmarshal(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
