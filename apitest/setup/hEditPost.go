package setup

import (
	"github.com/romshark/dgraph_graphql_go/api/graph"
	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/stretchr/testify/require"
)

func (h Helper) editPost(
	assumedSuccess successAssumption,
	postID store.ID,
	editorID store.ID,
	newTitle *string,
	newContents *string,
) (*gqlmod.Post, *graph.ResponseError) {
	t := h.c.t

	var old struct {
		Post *gqlmod.Post `json:"post"`
	}
	require.NoError(t, h.ts.Debug().QueryVar(
		`query($postId: Identifier!) {
			post(id: $postId) {
				id
				title
				contents
				creation
				author {
					id
				}
			}
		}`,
		map[string]interface{}{
			"postId": string(postID),
		},
		&old,
	))

	var result struct {
		EditPost *gqlmod.Post `json:"editPost"`
	}
	err := h.c.QueryVar(
		`mutation (
			$post: Identifier!
			$editor: Identifier!
			$newTitle: String
			$newContents: String
		) {
			editPost(
				post: $post
				editor: $editor
				newTitle: $newTitle
				newContents: $newContents
			) {
				id
				title
				contents
				creation
				author {
					id
				}
			}
		}`,
		map[string]interface{}{
			"post":        string(postID),
			"editor":      string(editorID),
			"newTitle":    newTitle,
			"newContents": newContents,
		},
		&result,
	)

	if err := checkErr(t, assumedSuccess, err); err != nil {
		return nil, err
	}

	require.NotNil(t, result.EditPost)
	if old.Post != nil {
		require.Equal(t, *old.Post.ID, *result.EditPost.ID)
		if newTitle != nil {
			require.Equal(t, *newTitle, *result.EditPost.Title)
		} else {
			require.Equal(t, *old.Post.Title, *result.EditPost.Title)
		}
		if newContents != nil {
			require.Equal(t, *newContents, *result.EditPost.Contents)
		} else {
			require.Equal(t, *old.Post.Contents, *result.EditPost.Contents)
		}
		require.Equal(t, *old.Post.Author.ID, *result.EditPost.Author.ID)
		require.Equal(t, *old.Post.Creation, *result.EditPost.Creation)
	}

	return result.EditPost, nil
}

// EditPost helps edit a post
func (h Helper) EditPost(
	postID store.ID,
	editorID store.ID,
	newTitle *string,
	newContents *string,
) (*gqlmod.Post, *graph.ResponseError) {
	return h.editPost(
		potentialFailure,
		postID,
		editorID,
		newTitle,
		newContents,
	)
}

// EditPost helps edit a post and assumes success
func (ok AssumeSuccess) EditPost(
	postID store.ID,
	editorID store.ID,
	newTitle *string,
	newContents *string,
) *gqlmod.Post {
	result, _ := ok.h.editPost(
		success,
		postID,
		editorID,
		newTitle,
		newContents,
	)
	return result
}
