package graphql

import "math/rand"

func rndInt64(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}
