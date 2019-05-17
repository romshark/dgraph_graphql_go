package setup

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
	"github.com/romshark/dgraph_graphql_go/store/errors"
	"github.com/stretchr/testify/require"
)

func (h Helper) createReaction(
	expectedErrorCode errors.Code,
	author store.ID,
	subject store.ID,
	emotion emotion.Emotion,
	message string,
) *gqlmod.Reaction {
	t := h.c.t

	var result struct {
		CreateReaction *gqlmod.Reaction `json:"createReaction"`
	}
	checkErr(t, expectedErrorCode, h.c.QueryVar(
		`mutation (
			$author: Identifier!
			$subject: Identifier!
			$emotion: Emotion!
			$message: String!
		) {
			createReaction(
				author: $author
				subject: $subject
				emotion: $emotion
				message: $message
			) {
				id
				creation
				message
				emotion
			}
		}`,
		map[string]interface{}{
			"author":  string(author),
			"subject": string(subject),
			"emotion": string(emotion),
			"message": message,
		},
		&result,
	))

	if expectedErrorCode != "" {
		return nil
	}

	require.NotNil(t, result.CreateReaction)
	require.Len(t, *result.CreateReaction.ID, 32)
	require.Equal(t, emotion, *result.CreateReaction.Emotion)
	require.Equal(t, message, *result.CreateReaction.Message)
	require.WithinDuration(
		t,
		time.Now(),
		*result.CreateReaction.Creation,
		h.creationTimeTollerance,
	)

	return result.CreateReaction
}

// CreateReaction helps creating a reaction and assumes success
func (ok AssumeSuccess) CreateReaction(
	author store.ID,
	subject store.ID,
	emotion emotion.Emotion,
	message string,
) *gqlmod.Reaction {
	return ok.h.createReaction(
		"",
		author,
		subject,
		emotion,
		message,
	)
}

// CreateReaction helps creating a reaction
func (notOk AssumeFailure) CreateReaction(
	expectedErrorCode errors.Code,
	author store.ID,
	subject store.ID,
	emotion emotion.Emotion,
	message string,
) {
	notOk.checkErrCode(expectedErrorCode)
	notOk.h.createReaction(
		expectedErrorCode,
		author,
		subject,
		emotion,
		message,
	)
}
