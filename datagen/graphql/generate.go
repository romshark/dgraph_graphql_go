package graphql

import (
	"context"
	"fmt"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/datagen"
)

func (gen *gqlgen) Generate(
	ctx context.Context,
	opts datagen.Options,
	stats datagen.StatisticsWriter,
) (graph datagen.Graph, err error) {
	// Prepare
	if stats == nil {
		stats = datagen.NewStatisticsRecorder()
	}

	graph.Users = make([]gqlmod.User, opts.Users)
	graph.UserPasswords = make([]string, opts.Users)
	graph.Posts = make([]gqlmod.Post, opts.Posts)

	if err = gen.generateUsers(ctx, &graph); err != nil {
		err = fmt.Errorf("user generation: %s", err)
		return
	}
	if err = gen.createUsers(ctx, &graph, stats); err != nil {
		err = fmt.Errorf("user creation: %s", err)
		return
	}
	if err = gen.generatePosts(ctx, &graph); err != nil {
		err = fmt.Errorf("post generation: %s", err)
		return
	}
	if err = gen.createPosts(ctx, &graph, stats); err != nil {
		err = fmt.Errorf("post creation: %s", err)
		return
	}

	// Print emtpy line to terminate the last progress bar line
	fmt.Println("")

	return
}
