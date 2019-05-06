package http

import (
	"time"
)

// ServerOptions defines the HTTP server transport layer options
type ServerOptions struct {
	Host              string
	KeepAliveDuration time.Duration
}

// SetDefaults sets the default options
func (opts *ServerOptions) SetDefaults() {
	if opts.KeepAliveDuration == time.Duration(0) {
		opts.KeepAliveDuration = 3 * time.Minute
	}
}

// ClientOptions defines the HTTP client transport layer options
type ClientOptions struct {
	Timeout time.Duration
}

// SetDefaults sets the default options
func (opts *ClientOptions) SetDefaults() {
	if opts.Timeout == time.Duration(0) {
		opts.Timeout = 30 * time.Second
	}
}
