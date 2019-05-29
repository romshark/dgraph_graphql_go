package config

import (
	"crypto/tls"
	"fmt"
	"reflect"
)

// TLSCurveID represents a TLS curve identifier
type TLSCurveID tls.CurveID

// UnmarshalTOML implements the TOML unmarshaler interface
func (v *TLSCurveID) UnmarshalTOML(val interface{}) error {
	if str, isString := val.(string); isString {
		switch str {
		case "CurveP256":
			*v = TLSCurveID(tls.CurveP256)
		case "CurveP384":
			*v = TLSCurveID(tls.CurveP384)
		case "CurveP521":
			*v = TLSCurveID(tls.CurveP521)
		case "X25519":
			*v = TLSCurveID(tls.X25519)
		default:
			return fmt.Errorf("unknown TLS curve ID: '%s'", val)
		}
		return nil
	}
	return fmt.Errorf(
		"unexpected TLS curve value type: %s",
		reflect.TypeOf(val),
	)
}
