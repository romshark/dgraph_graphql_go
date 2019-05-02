package sesskeygen

// SessionKeyGenerator defines the session key generator interface
type SessionKeyGenerator interface {
	// Generates a new cryptographically safe session key
	Generate() string
}
