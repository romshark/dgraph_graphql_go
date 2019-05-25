package graphql

import (
	"context"
	"log"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/datagen"
	"github.com/schollz/progressbar"
)

func (gen *gqlgen) generateUsers(
	ctx context.Context,
	graph *datagen.Graph,
) error {
	// Generate user profiles
	log.Println("generate users")
	progressUserGen := progressbar.New(len(graph.Users))
	progressUserGen.RenderBlank()
	for i := range graph.Users {
		email := gen.emailGen.New()
		displayName := gen.displayNameGen.New()
		graph.Users[i] = gqlmod.User{
			Email:       &email,
			DisplayName: &displayName,
		}
		graph.UserPasswords[i] = "testpass"
		progressUserGen.Add(1)
	}

	return nil
}
