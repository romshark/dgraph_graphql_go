package setup

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/errors"
	"github.com/stretchr/testify/require"
)

func (h Helper) createPost(
	expectedErrorCode errors.Code,
	authorID store.ID,
	title string,
	contents string,
) *gqlmod.Post {
	t := h.c.t

	var result struct {
		CreatePost *gqlmod.Post `json:"createPost"`
	}
	checkErr(t, expectedErrorCode, h.c.QueryVar(
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
		map[string]interface{}{
			"author":   string(authorID),
			"title":    title,
			"contents": contents,
		},
		&result,
	))

	if expectedErrorCode != "" {
		return nil
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

	return result.CreatePost
}

// CreatePost helps creating a user and assumes success
func (ok AssumeSuccess) CreatePost(
	authorID store.ID,
	title string,
	contents string,
) *gqlmod.Post {
	return ok.h.createPost("", authorID, title, contents)
}

// CreatePost assumes the given error code to be returned
func (notOk AssumeFailure) CreatePost(
	expectedErrorCode errors.Code,
	authorID store.ID,
	title string,
	contents string,
) {
	notOk.checkErrCode(expectedErrorCode)
	notOk.h.createPost(expectedErrorCode, authorID, title, contents)
}
