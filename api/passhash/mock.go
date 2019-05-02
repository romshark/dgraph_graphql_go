package passhash

import "bytes"

// Mock implements the PasswordHasher interface using bcrypt
type Mock struct{}

// Hash salts and hashes the given password returning the resulting hash
func (h Mock) Hash(password []byte) ([]byte, error) {
	return password, nil
}

// Compare returns true if the given password corresponds to the given hash,
// otherwise returns false
func (h Mock) Compare(password, hash []byte) bool {
	return bytes.Equal(password, hash)
}
