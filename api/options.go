package api

import (
	"errors"

	"github.com/romshark/dgraph_graphql_go/api/passhash"
	"github.com/romshark/dgraph_graphql_go/api/sesskeygen"
	"github.com/romshark/dgraph_graphql_go/api/transport"
	thttp "github.com/romshark/dgraph_graphql_go/api/transport/http"
)

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

// DebugUserStatus defines the debug user status option
type DebugUserStatus string

const (
	// DebugUserUnset represents the default unset option value
	DebugUserUnset DebugUserStatus = ""

	// DebugUserDisabled disables the debug user
	DebugUserDisabled DebugUserStatus = "disabled"

	// DebugUserReadOnly enables the debug user in a read-only mode
	DebugUserReadOnly DebugUserStatus = "read-only"

	// DebugUserRW enables the debug user in a read-write mode
	DebugUserRW DebugUserStatus = "read-write"
)

// DebugUserOptions defines the API debug user options
type DebugUserOptions struct {
	Status   DebugUserStatus
	Username string
	Password string
}

// Prepares sets defaults and validates the options
func (opts *DebugUserOptions) Prepares(mode Mode) error {
	// Set default debug user option
	if opts.Status == DebugUserUnset {
		switch mode {
		case ModeProduction:
			opts.Status = DebugUserDisabled
		case ModeBeta:
			opts.Status = DebugUserReadOnly
		default:
			opts.Status = DebugUserRW
		}
	}

	// Use "debug" as the default debug username
	if opts.Username == "" {
		opts.Username = "debug"
	}

	// Use "debug" as the default debug password
	if opts.Password == "" {
		opts.Password = "debug"
	}

	// VALIDATE

	// Ensure the debug user isn't enabled in production mode
	if mode == ModeProduction {
		if opts.Status != DebugUserDisabled {
			return errors.New("debug user must be disabled in production mode")
		}
	}

	return nil
}

// ServerOptions defines the API server options
type ServerOptions struct {
	Mode                Mode
	Host                string
	DBHost              string
	SessionKeyGenerator sesskeygen.SessionKeyGenerator
	PasswordHasher      passhash.PasswordHasher
	DebugUser           DebugUserOptions
	Transport           []transport.Server
}

// Prepare sets defaults and validates the options
func (opts *ServerOptions) Prepare() error {
	// Use production mode by default
	if opts.Mode == "" {
		opts.Mode = ModeProduction
	}

	// Set default database host address
	if opts.DBHost == "" {
		switch opts.Mode {
		case ModeProduction:
			opts.DBHost = "localhost:8080"
		default:
			opts.DBHost = "localhost:10080"
		}
	}

	// Use default session key generator
	if opts.SessionKeyGenerator == nil {
		opts.SessionKeyGenerator = sesskeygen.NewDefault()
	}

	// Use default password hasher
	if opts.PasswordHasher == nil {
		opts.PasswordHasher = passhash.Bcrypt{}
	}

	// VALIDATE

	// Ensure at least one transport adapter is specified
	if len(opts.Transport) < 1 {
		return errors.New("no transport adapter")
	}

	if opts.Mode == ModeProduction {
		// Validate transport adapters
		for _, trn := range opts.Transport {
			if httpAdapter, ok := trn.(*thttp.Server); ok {
				// Ensure TLS is enabled in production on all transport adapters
				opts := httpAdapter.Options()
				if opts.TLS == nil {
					return errors.New(
						"TLS must not be disabled on HTTP transport adapter " +
							"in production mode",
					)
				}

				// Ensure playground is disabled in production
				if opts.Playground {
					return errors.New(
						"the playground must be disabled on " +
							"HTTP transport adapter in production mode",
					)
				}
			}
		}

		// Ensure standard session key generator is used in production
		if _, ok := opts.SessionKeyGenerator.(*sesskeygen.Default); !ok {
			return errors.New(
				"standard session key generator " +
					"must be used in production mode",
			)
		}

		// Ensure bcrypt password hasher is used in production
		if _, ok := opts.PasswordHasher.(passhash.Bcrypt); !ok {
			return errors.New(
				"bcrypt password hasher must be used in production mode",
			)
		}
	}

	return opts.DebugUser.Prepares(opts.Mode)
}
