package config

import (
	"errors"
	"log"
	"os"

	"github.com/romshark/dgraph_graphql_go/api/passhash"
	"github.com/romshark/dgraph_graphql_go/api/sesskeygen"
	"github.com/romshark/dgraph_graphql_go/api/transport"
	thttp "github.com/romshark/dgraph_graphql_go/api/transport/http"
)

// ServerConfig defines the API server configurations
type ServerConfig struct {
	Mode                Mode
	DBHost              string
	SessionKeyGenerator sesskeygen.SessionKeyGenerator
	PasswordHasher      passhash.PasswordHasher
	DebugUser           DebugUserConfig
	Transport           []transport.Server
	DebugLog            *log.Logger
	ErrorLog            *log.Logger
}

// Prepare sets defaults and validates the configurations
func (conf *ServerConfig) Prepare() error {
	// Use production mode by default
	if conf.Mode == "" {
		conf.Mode = ModeProduction
	}

	// Set default database host address
	if conf.DBHost == "" {
		switch conf.Mode {
		case ModeProduction:
			conf.DBHost = "localhost:9080"
		default:
			conf.DBHost = "localhost:10180"
		}
	}

	// Use default session key generator
	if conf.SessionKeyGenerator == nil {
		conf.SessionKeyGenerator = sesskeygen.NewDefault()
	}

	// Use default password hasher
	if conf.PasswordHasher == nil {
		conf.PasswordHasher = passhash.Bcrypt{}
	}

	// Use default debug logger to stdout
	if conf.DebugLog == nil {
		conf.DebugLog = log.New(
			os.Stdout,
			"DBG: ",
			log.Ldate|log.Ltime,
		)
	}

	// Use default error logger to stderr
	if conf.ErrorLog == nil {
		conf.ErrorLog = log.New(
			os.Stderr,
			"ERR: ",
			log.Ldate|log.Ltime,
		)
	}

	// VALIDATE

	// Ensure at least one transport adapter is specified
	if len(conf.Transport) < 1 {
		return errors.New("no transport adapter")
	}

	if conf.Mode == ModeProduction {
		// Validate transport adapters
		for _, trn := range conf.Transport {
			if httpAdapter, ok := trn.(*thttp.Server); ok {
				// Ensure TLS is enabled in production on all transport adapters
				conf := httpAdapter.Config()
				if conf.TLS == nil {
					return errors.New(
						"TLS must not be disabled on HTTP transport adapter " +
							"in production mode",
					)
				}

				// Ensure playground is disabled in production
				if conf.Playground {
					return errors.New(
						"the playground must be disabled on " +
							"HTTP transport adapter in production mode",
					)
				}
			}
		}

		// Ensure standard session key generator is used in production
		if _, ok := conf.SessionKeyGenerator.(*sesskeygen.Default); !ok {
			return errors.New(
				"standard session key generator " +
					"must be used in production mode",
			)
		}

		// Ensure bcrypt password hasher is used in production
		if _, ok := conf.PasswordHasher.(passhash.Bcrypt); !ok {
			return errors.New(
				"bcrypt password hasher must be used in production mode",
			)
		}
	}

	return conf.DebugUser.Prepares(conf.Mode)
}
