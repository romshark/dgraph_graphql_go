package http

import (
	"crypto/tls"
	"errors"
	"time"
)

// ServerTLS represents the TLS configurations
type ServerTLS struct {
	Config              *tls.Config
	CertificateFilePath string
	PrivateKeyFilePath  string
}

// Clone creates an exact detached copy of the server TLS configurations
func (stls *ServerTLS) Clone() *ServerTLS {
	if stls == nil {
		return nil
	}
	var config *tls.Config
	if stls.Config != nil {
		config = stls.Config.Clone()
	}
	return &ServerTLS{
		Config:              config,
		CertificateFilePath: stls.CertificateFilePath,
		PrivateKeyFilePath:  stls.PrivateKeyFilePath,
	}
}

// ServerConfig defines the HTTP server transport layer configurations
type ServerConfig struct {
	Host              string
	KeepAliveDuration time.Duration
	TLS               *ServerTLS
	Playground        bool
}

// Prepare sets defaults and validates the configurations
func (conf *ServerConfig) Prepare() error {
	if conf.KeepAliveDuration == time.Duration(0) {
		conf.KeepAliveDuration = 3 * time.Minute
	}

	if conf.TLS != nil {
		if conf.TLS.CertificateFilePath == "" {
			return errors.New("missing TLS certificate file path")
		}
		if conf.TLS.PrivateKeyFilePath == "" {
			return errors.New("missing TLS private key file path")
		}
	}

	return nil
}

// ClientConfig defines the HTTP client transport layer configuration
type ClientConfig struct {
	Timeout time.Duration
}

// SetDefaults sets the default configuration
func (conf *ClientConfig) SetDefaults() {
	if conf.Timeout == time.Duration(0) {
		conf.Timeout = 30 * time.Second
	}
}
