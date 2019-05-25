package graphql

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/icrowley/fake"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type postContentsGenerator struct {
	maxLen uint32
}

func newPostContentsGenerator(maxLen uint32) *postContentsGenerator {
	if maxLen < 1 {
		panic(fmt.Errorf("invalid maxLen: %d", maxLen))
	}
	return &postContentsGenerator{
		maxLen: maxLen,
	}
}

func (gen *postContentsGenerator) New() (newPostContents string) {
	var r string
	length := rndInt64(int64(gen.maxLen/2), int64(gen.maxLen))
	for {
		r += fake.Word() + " "
		if int64(len(r)) > length {
			return r[:length]
		}
	}
}
