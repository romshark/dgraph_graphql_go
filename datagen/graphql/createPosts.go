package graphql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/datagen"
	"github.com/schollz/progressbar"
)

func (gen *gqlgen) createPosts(
	ctx context.Context,
	graph *datagen.Graph,
	stats datagen.StatisticsWriter,
) error {
	wg := NewErrWaitGroup(uint64(len(graph.Posts)))

	// Create posts
	fmt.Println("")
	log.Println("create posts")
	progressPostCrt := progressbar.New(len(graph.Posts))
	progressPostCrt.RenderBlank()

	for i, p := range graph.Posts {
		index := i
		post := p
		go func() {
			gen.slots.Acquire(ctx, 1)
			defer gen.slots.Release(1)

			start := time.Now()
			var qr struct {
				NewPost gqlmod.Post `json:"createPost"`
			}
			if err := gen.client.QueryVar(
				`mutation(
					$author: Identifier!
					$title: String!
					$contents: String!
				) {
					createPost(
						author: $author,
						title: $title,
						contents: $contents
					) {
						id
						title
						contents
						creation
					}
				}`,
				map[string]interface{}{
					"author":   *post.Author.ID,
					"title":    *post.Title,
					"contents": *post.Contents,
				},
				&qr,
			); err != nil {
				wg.Fail(err)
				return
			}
			graph.Posts[index] = qr.NewPost
			post.Author.Posts = append(post.Author.Posts, qr.NewPost)
			stats.RecordPostCreation(time.Since(start))
			progressPostCrt.Add(1)
			wg.Dec(1)
		}()
	}

	return wg.Wait()
}
