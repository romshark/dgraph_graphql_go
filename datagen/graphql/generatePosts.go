package graphql

import (
	"context"
	"fmt"
	"log"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/datagen"
	"github.com/schollz/progressbar"
)

func (gen *gqlgen) generatePosts(
	ctx context.Context,
	graph *datagen.Graph,
) error {
	// Generate posts
	fmt.Println("")
	log.Println("generate posts")
	progressPostGen := progressbar.New(len(graph.Users))
	progressPostGen.RenderBlank()
	for i := range graph.Posts {
		title := gen.postTitleGen.New()
		contents := gen.postContentsGen.New()
		graph.Posts[i] = gqlmod.Post{
			Title:    &title,
			Contents: &contents,
			Author:   graph.GetRndUser(),
		}
		progressPostGen.Add(1)
	}

	return nil
}
