package api

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/api/passhash"
	"github.com/romshark/dgraph_graphql_go/api/sesskeygen"
)

// RootUserStatus defines the root user status option
type RootUserStatus byte

const (
	// RootUserUnset represents the default unset option value
	RootUserUnset RootUserStatus = 0

	// RootUserDisabled disables the root user
	RootUserDisabled

	// RootUserReadOnly enables the root user in a read-only mode
	RootUserReadOnly RootUserStatus = 2

	// RootUserRW enables the root user in a read-write mode
	RootUserRW RootUserStatus = 3
)

// RootUserOptions defines the API root user options
type RootUserOptions struct {
	Status   RootUserStatus
	Username string
	Password string
}

// SetDefaults sets the default options
func (opts *RootUserOptions) SetDefaults() {
	// Disable the root user by default
	if opts.Status == RootUserUnset {
		opts.Status = RootUserDisabled
	}

	// Use "root" as the default root username
	if opts.Username == "" {
		opts.Username = "root"
	}

	// Use "root" as the default root password
	if opts.Password == "" {
		opts.Password = "root"
	}
}

// ServerOptions defines the API server options
type ServerOptions struct {
	Host                string
	DBHost              string
	SessionKeyGenerator sesskeygen.SessionKeyGenerator
	PasswordHasher      passhash.PasswordHasher
	KeepAliveDuration   time.Duration
	RootUser            RootUserOptions
}

// SetDefaults sets the default options
func (opts *ServerOptions) SetDefaults() {
	if opts.KeepAliveDuration == time.Duration(0) {
		opts.KeepAliveDuration = 3 * time.Minute
	}

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
}
