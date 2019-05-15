package api

import (
	"github.com/romshark/dgraph_graphql_go/api/passhash"
	"github.com/romshark/dgraph_graphql_go/api/sesskeygen"
	"github.com/romshark/dgraph_graphql_go/api/transport"
	thttp "github.com/romshark/dgraph_graphql_go/api/transport/http"
)

// DebugUserStatus defines the debug user status option
type DebugUserStatus byte

const (
	// DebugUserUnset represents the default unset option value
	DebugUserUnset DebugUserStatus = 0

	// DebugUserDisabled disables the debug user
	DebugUserDisabled

	// DebugUserReadOnly enables the debug user in a read-only mode
	DebugUserReadOnly DebugUserStatus = 2

	// DebugUserRW enables the debug user in a read-write mode
	DebugUserRW DebugUserStatus = 3
)

// DebugUserOptions defines the API debug user options
type DebugUserOptions struct {
	Status   DebugUserStatus
	Username string
	Password string
}

// SetDefaults sets the default options
func (opts *DebugUserOptions) SetDefaults() {
	// Disable the debug user by default
	if opts.Status == DebugUserUnset {
		opts.Status = DebugUserDisabled
	}

	// Use "debug" as the default debug username
	if opts.Username == "" {
		opts.Username = "debug"
	}

	// Use "debug" as the default debug password
	if opts.Password == "" {
		opts.Password = "debug"
	}
}

// ServerOptions defines the API server options
type ServerOptions struct {
	Host                string
	DBHost              string
	SessionKeyGenerator sesskeygen.SessionKeyGenerator
	PasswordHasher      passhash.PasswordHasher
	DebugUser           DebugUserOptions
	Transport           []transport.Server
}

// SetDefaults sets the default options
func (opts *ServerOptions) SetDefaults() error {
	// Use default non-production database port
	if opts.DBHost == "" {
		opts.DBHost = "localhost:6000"
	}

	// Use default session key generator
	if opts.SessionKeyGenerator == nil {
		opts.SessionKeyGenerator = sesskeygen.NewDefault()
	}

	// Use default password hasher
	if opts.PasswordHasher == nil {
		opts.PasswordHasher = passhash.Bcrypt{}
	}

	// Use HTTP as the default transport
	if len(opts.Transport) < 1 {
		httpTransport, err := thttp.NewServer(thttp.ServerOptions{})
		if err != nil {
			return err
		}
		opts.Transport = []transport.Server{
			httpTransport,
		}
	}

	return nil
}
