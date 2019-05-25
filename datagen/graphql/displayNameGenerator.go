package graphql

import (
	"github.com/Pallinder/go-randomdata"
)

type displayNameGenerator struct {
	displayNames map[string]struct{}
}

func newDisplayNameGenerator() *displayNameGenerator {
	return &displayNameGenerator{
		displayNames: make(map[string]struct{}),
	}
}

func (gen *displayNameGenerator) New() (newDisplayName string) {
	for {
		newDisplayName = randomdata.SillyName()
		if _, isIn := gen.displayNames[newDisplayName]; !isIn {
			gen.displayNames[newDisplayName] = struct{}{}
			break
		}
	}
	return
}
