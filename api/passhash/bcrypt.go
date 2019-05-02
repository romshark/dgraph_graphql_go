package passhash

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Bcrypt implements the PasswordHasher interface using bcrypt
type Bcrypt struct{}

// Hash salts and hashes the given password returning the resulting hash
func (h Bcrypt) Hash(password []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, errors.Wrap(err, "generate password hash")
	}
	return hash, nil
}

// Compare returns true if the given password corresponds to the given hash,
// otherwise returns false
func (h Bcrypt) Compare(password, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}
