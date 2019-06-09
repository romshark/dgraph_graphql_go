package config

import (
	"crypto/tls"
	"fmt"
	"reflect"
)

// TLSCipherSuite represents a TLS cipher suite
type TLSCipherSuite uint16

// UnmarshalTOML implements the TOML unmarshaler interface
func (v *TLSCipherSuite) UnmarshalTOML(val interface{}) error {
	if str, isString := val.(string); isString {
		switch str {
		case "RSA_WITH_RC4_128_SHA":
			*v = TLSCipherSuite(tls.TLS_RSA_WITH_RC4_128_SHA)
		case "RSA_WITH_3DES_EDE_CBC_SHA":
			*v = TLSCipherSuite(tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA)
		case "RSA_WITH_AES_128_CBC_SHA":
			*v = TLSCipherSuite(tls.TLS_RSA_WITH_AES_128_CBC_SHA)
		case "RSA_WITH_AES_256_CBC_SHA":
			*v = TLSCipherSuite(tls.TLS_RSA_WITH_AES_256_CBC_SHA)
		case "RSA_WITH_AES_128_CBC_SHA256":
			*v = TLSCipherSuite(tls.TLS_RSA_WITH_AES_128_CBC_SHA256)
		case "RSA_WITH_AES_128_GCM_SHA256":
			*v = TLSCipherSuite(tls.TLS_RSA_WITH_AES_128_GCM_SHA256)
		case "RSA_WITH_AES_256_GCM_SHA384":
			*v = TLSCipherSuite(tls.TLS_RSA_WITH_AES_256_GCM_SHA384)
		case "ECDHE_ECDSA_WITH_RC4_128_SHA":
			*v = TLSCipherSuite(tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA)
		case "ECDHE_ECDSA_WITH_AES_128_CBC_SHA":
			*v = TLSCipherSuite(tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA)
		case "ECDHE_ECDSA_WITH_AES_256_CBC_SHA":
			*v = TLSCipherSuite(tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA)
		case "ECDHE_RSA_WITH_RC4_128_SHA":
			*v = TLSCipherSuite(tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA)
		case "ECDHE_RSA_WITH_3DES_EDE_CBC_SHA":
			*v = TLSCipherSuite(tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA)
		case "ECDHE_RSA_WITH_AES_128_CBC_SHA":
			*v = TLSCipherSuite(tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA)
		case "ECDHE_RSA_WITH_AES_256_CBC_SHA":
			*v = TLSCipherSuite(tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA)
		case "ECDHE_ECDSA_WITH_AES_128_CBC_SHA256":
			*v = TLSCipherSuite(tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256)
		case "ECDHE_RSA_WITH_AES_128_CBC_SHA256":
			*v = TLSCipherSuite(tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256)
		case "ECDHE_RSA_WITH_AES_128_GCM_SHA256":
			*v = TLSCipherSuite(tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256)
		case "ECDHE_ECDSA_WITH_AES_128_GCM_SHA256":
			*v = TLSCipherSuite(tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256)
		case "ECDHE_RSA_WITH_AES_256_GCM_SHA384":
			*v = TLSCipherSuite(tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384)
		case "ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":
			*v = TLSCipherSuite(tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384)
		case "ECDHE_RSA_WITH_CHACHA20_POLY1305":
			*v = TLSCipherSuite(tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305)
		case "ECDHE_ECDSA_WITH_CHACHA20_POLY1305":
			*v = TLSCipherSuite(tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305)
		case "AES_128_GCM_SHA256":
			*v = TLSCipherSuite(tls.TLS_AES_128_GCM_SHA256)
		case "AES_256_GCM_SHA384":
			*v = TLSCipherSuite(tls.TLS_AES_256_GCM_SHA384)
		case "CHACHA20_POLY1305_SHA256":
			*v = TLSCipherSuite(tls.TLS_CHACHA20_POLY1305_SHA256)
		default:
			return fmt.Errorf("unknown TLS cipher suite: '%s'", val)
		}
		return nil
	}
	return fmt.Errorf(
		"unexpected TLS cipher suite value type: %s",
		reflect.TypeOf(val),
	)
}
