package config

import (
	"crypto/tls"
	"fmt"
	"reflect"
)

// TLSVersion represents a TLS protocol version
type TLSVersion uint16

// UnmarshalTOML implements the TOML unmarshaler interface
func (v *TLSVersion) UnmarshalTOML(val interface{}) error {
	if str, isString := val.(string); isString {
		switch str {
		case "SSL 3.0":
			*v = TLSVersion(tls.VersionSSL30)
		case "TLS 1.0":
			*v = TLSVersion(tls.VersionTLS10)
		case "TLS 1.1":
			*v = TLSVersion(tls.VersionTLS11)
		case "TLS 1.2":
			*v = TLSVersion(tls.VersionTLS12)
		case "TLS 1.3":
			*v = TLSVersion(tls.VersionTLS13)
		default:
			return fmt.Errorf("unknown TLS protocol version: '%s'", val)
		}
		return nil
	}
	return fmt.Errorf(
		"unexpected TLS protocol version value type: %s",
		reflect.TypeOf(val),
	)
}
