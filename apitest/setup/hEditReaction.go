package setup

import (
	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/errors"
	"github.com/stretchr/testify/require"
)

func (h Helper) editReaction(
	expectedErrorCode errors.Code,
	reactionID store.ID,
	editorID store.ID,
	newMessage string,
) *gqlmod.Reaction {
	t := h.c.t

	var old struct {
		Reaction *gqlmod.Reaction `json:"reaction"`
	}
	require.NoError(t, h.ts.Debug().QueryVar(
		`query($reactionID: Identifier!) {
			reaction(id: $reactionID) {
				id
				emotion
				message
				creation
				author {
					id
				}
				subject {
					__typename
					... on Post {
						id
					}
					... on Reaction {
						id
					}
				}
			}
		}`,
		map[string]interface{}{
			"reactionID": string(reactionID),
		},
		&old,
	))

	var result struct {
		EditReaction *gqlmod.Reaction `json:"editReaction"`
	}
	checkErr(t, expectedErrorCode, h.c.QueryVar(
		`mutation (
			$reaction: Identifier!
			$editor: Identifier!
			$newMessage: String!
		) {
			editReaction(
				reaction: $reaction
				editor: $editor
				newMessage: $newMessage
			) {
				id
				emotion
				subject {
					__typename
					... on Post {
						id
					}
					... on Reaction {
						id
					}
				}
				message
				creation
				author {
					id
				}
			}
		}`,
		map[string]interface{}{
			"reaction":   string(reactionID),
			"editor":     string(editorID),
			"newMessage": newMessage,
		},
		&result,
	))

	if expectedErrorCode != "" {
		return nil
	}

	require.NotNil(t, result.EditReaction)
	if old.Reaction != nil {
		require.Equal(t, *old.Reaction.ID, *result.EditReaction.ID)
		require.Equal(t, *old.Reaction.Emotion, *result.EditReaction.Emotion)
		require.Equal(t, old.Reaction.Subject, result.EditReaction.Subject)
		require.Equal(t, newMessage, *result.EditReaction.Message)
		require.Equal(
			t,
			*old.Reaction.Author.ID,
			*result.EditReaction.Author.ID,
		)
		require.Equal(t, *old.Reaction.Creation, *result.EditReaction.Creation)
	}

	return result.EditReaction
}

// EditReaction helps edit a reaction and assumes success
func (ok AssumeSuccess) EditReaction(
	reactionID store.ID,
	editorID store.ID,
	newMessage string,
) *gqlmod.Reaction {
	return ok.h.editReaction(
		"",
		reactionID,
		editorID,
		newMessage,
	)
}

// EditReaction helps edit a reaction
func (notOk AssumeFailure) EditReaction(
	expectedErrorCode errors.Code,
	reactionID store.ID,
	editorID store.ID,
	newMessage string,
) {
	notOk.checkErrCode(expectedErrorCode)
	notOk.h.editReaction(
		expectedErrorCode,
		reactionID,
		editorID,
		newMessage,
	)
}
