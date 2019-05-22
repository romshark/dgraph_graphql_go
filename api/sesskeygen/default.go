package sesskeygen

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"fmt"
)

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(length uint32) (bytes []byte, err error) {
	bytes = make([]byte, length)
	_, err = cryptoRand.Read(bytes)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// generateSessionKey returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateSessionKey() string {
	bytes, err := generateRandomBytes(48)
	if err != nil {
		panic(fmt.Errorf("could not generate a session key"))
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// Default represents a default implementation of SessionKeyGenerator
type Default struct{}

// NewDefault creates a new default session key generator instance
func NewDefault() SessionKeyGenerator {
	return &Default{}
}

// Generate implements the SessionKeyGenerator interface
func (gen *Default) Generate() string {
	return generateSessionKey()
}
