package http

import (
	"crypto/tls"
	"errors"
	"time"
)

// ServerTLS represents the TLS options
type ServerTLS struct {
	Config              *tls.Config
	CertificateFilePath string
	PrivateKeyFilePath  string
}

// Clone creates an exact detached copy of the server TLS options
func (stls *ServerTLS) Clone() *ServerTLS {
	if stls == nil {
		return nil
	}
	return &ServerTLS{
		Config:              stls.Config.Clone(),
		CertificateFilePath: stls.CertificateFilePath,
		PrivateKeyFilePath:  stls.PrivateKeyFilePath,
	}
}

// ServerOptions defines the HTTP server transport layer options
type ServerOptions struct {
	Host              string
	KeepAliveDuration time.Duration
	TLS               *ServerTLS
	Playground        bool
}

// Prepare sets defaults and validates the options
func (opts *ServerOptions) Prepare() error {
	if opts.KeepAliveDuration == time.Duration(0) {
		opts.KeepAliveDuration = 3 * time.Minute
	}

	if opts.TLS != nil {
		if opts.TLS.CertificateFilePath == "" {
			return errors.New("missing TLS certificate file path")
		}
		if opts.TLS.PrivateKeyFilePath == "" {
			return errors.New("missing TLS private key file path")
		}
	}

	return nil
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
