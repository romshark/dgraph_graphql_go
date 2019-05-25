package datagen

import "github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"

// Graph represents the result of a successful generation
type Graph struct {
	Users         []gqlmod.User
	UserPasswords []string
	Posts         []gqlmod.Post
}

// GetRndUser returns a random generated user
func (gr *Graph) GetRndUser() *gqlmod.User {
	return &gr.Users[rndInt64(0, int64(len(gr.Users)))]
}

// GetRndPost returns a random generated post
func (gr *Graph) GetRndPost() *gqlmod.Post {
	return &gr.Posts[rndInt64(0, int64(len(gr.Posts)))]
}
