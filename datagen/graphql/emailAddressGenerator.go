package graphql

import (
	"github.com/Pallinder/go-randomdata"
)

type emailAddressGenerator struct {
	emails map[string]struct{}
}

func newEmailAddressGenerator() *emailAddressGenerator {
	return &emailAddressGenerator{
		emails: make(map[string]struct{}),
	}
}

func (gen *emailAddressGenerator) New() (newEmail string) {
	for {
		newEmail = randomdata.Email()
		if _, isIn := gen.emails[newEmail]; !isIn {
			gen.emails[newEmail] = struct{}{}
			break
		}
	}
	return
}
