package passhash

// PasswordHasher defines the interface of a password hasher
type PasswordHasher interface {
	// Hash must salt and hash the given password returning the resulting hash
	Hash(password []byte) ([]byte, error)

	// Compare must return true if the given password corresponds to
	// the given hash, otherwise must return false
	Compare(password, hash []byte) bool
}
