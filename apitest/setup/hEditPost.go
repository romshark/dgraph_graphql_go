package setup

import (
	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/errors"
	"github.com/stretchr/testify/require"
)

func (h Helper) editPost(
	expectedErrorCode errors.Code,
	postID store.ID,
	editorID store.ID,
	newTitle *string,
	newContents *string,
) *gqlmod.Post {
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
	checkErr(t, expectedErrorCode, h.c.QueryVar(
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
	))

	if expectedErrorCode != "" {
		return nil
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

	return result.EditPost
}

// EditPost helps edit a post and assumes success
func (ok AssumeSuccess) EditPost(
	postID store.ID,
	editorID store.ID,
	newTitle *string,
	newContents *string,
) *gqlmod.Post {
	return ok.h.editPost(
		"",
		postID,
		editorID,
		newTitle,
		newContents,
	)
}

// EditPost helps edit a post
func (notOk AssumeFailure) EditPost(
	expectedErrorCode errors.Code,
	postID store.ID,
	editorID store.ID,
	newTitle *string,
	newContents *string,
) {
	notOk.checkErrCode(expectedErrorCode)
	notOk.h.editPost(
		expectedErrorCode,
		postID,
		editorID,
		newTitle,
		newContents,
	)
}
