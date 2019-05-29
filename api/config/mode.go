package config

import "fmt"

// Mode defines the server mode
type Mode string

const (
	// ModeDebug represents the debug server mode
	ModeDebug Mode = "debug"

	// ModeBeta represents the beta server mode
	ModeBeta Mode = "beta"

	// ModeProduction represents the production server mode
	ModeProduction Mode = "production"
)

// Validate returns an error if the mode is unknown
func (mode Mode) Validate() error {
	switch mode {
	case ModeDebug:
		fallthrough
	case ModeBeta:
		fallthrough
	case ModeProduction:
		return nil
	}
	return fmt.Errorf("unknown mode: '%s'", mode)
}
