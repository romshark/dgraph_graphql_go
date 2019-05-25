package graphql

import (
	"github.com/Pallinder/go-randomdata"
)

type postTitleGenerator struct {
	postTitles map[string]struct{}
}

func newPostTitleGenerator() *postTitleGenerator {
	return &postTitleGenerator{
		postTitles: make(map[string]struct{}),
	}
}

func (gen *postTitleGenerator) New() (newPostTitle string) {
	for {
		newPostTitle = randomdata.SillyName()
		if _, isIn := gen.postTitles[newPostTitle]; !isIn {
			gen.postTitles[newPostTitle] = struct{}{}
			break
		}
	}
	return
}
