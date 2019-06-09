package config

import "errors"

// DebugUserConfig defines the API debug user configurations
type DebugUserConfig struct {
	Mode     DebugUserMode
	Username string
	Password string
}

// Prepares sets defaults and validates the configurations
func (conf *DebugUserConfig) Prepares(mode Mode) error {
	// Set default debug user mode
	if conf.Mode == DebugUserUnset {
		switch mode {
		case ModeProduction:
			conf.Mode = DebugUserDisabled
		case ModeBeta:
			conf.Mode = DebugUserReadOnly
		default:
			conf.Mode = DebugUserRW
		}
	}

	// Use "debug" as the default debug username
	if conf.Username == "" {
		conf.Username = "debug"
	}

	// Use "debug" as the default debug password
	if conf.Password == "" {
		conf.Password = "debug"
	}

	// VALIDATE

	// Ensure the debug user isn't enabled in production mode
	if mode == ModeProduction {
		if conf.Mode != DebugUserDisabled {
			return errors.New("debug user must be disabled in production mode")
		}
	}

	return nil
}
