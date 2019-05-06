package setup

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph"
	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/stretchr/testify/require"
)

func (h Helper) createPost(
	assumedSuccess successAssumption,
	authorID store.ID,
	title string,
	contents string,
) (*gqlmod.Post, *graph.ResponseError) {
	t := h.c.t

	var result struct {
		CreatePost *gqlmod.Post `json:"createPost"`
	}
	err := h.c.QueryVar(
		`mutation (
			$author: Identifier!
			$title: String!
			$contents: String!
		) {
			createPost(
				author: $author
				title: $title
				contents: $contents
			) {
				id
				title
				contents
				creation
				author {
					id
				}
				reactions {
					id
				}
			}
		}`,
		map[string]string{
			"author":   string(authorID),
			"title":    title,
			"contents": contents,
		},
		&result,
	)

	if err := checkErr(t, assumedSuccess, err); err != nil {
		return nil, err
	}

	require.NotNil(t, result.CreatePost)
	require.Len(t, *result.CreatePost.ID, 32)
	require.Equal(t, title, *result.CreatePost.Title)
	require.Equal(t, contents, *result.CreatePost.Contents)
	require.Equal(t, authorID, *result.CreatePost.Author.ID)
	require.Len(t, result.CreatePost.Reactions, 0)
	require.WithinDuration(
		t,
		time.Now(),
		*result.CreatePost.Creation,
		h.creationTimeTollerance,
	)

	return result.CreatePost, nil
}

// CreatePost helps creating a user
func (h Helper) CreatePost(
	authorID store.ID,
	title string,
	contents string,
) (*gqlmod.Post, *graph.ResponseError) {
	return h.createPost(potentialFailure, authorID, title, contents)
}

// CreatePost helps creating a user and assumes success
func (ok AssumeSuccess) CreatePost(
	authorID store.ID,
	title string,
	contents string,
) *gqlmod.Post {
	result, _ := ok.h.createPost(success, authorID, title, contents)
	return result
}
