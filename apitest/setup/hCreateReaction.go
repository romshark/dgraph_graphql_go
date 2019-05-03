package setup

import (
	"time"

	"github.com/romshark/dgraph_graphql_go/api"
	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
	"github.com/stretchr/testify/require"
)

func (h Helper) createReaction(
	successAssumption successAssumption,
	author store.ID,
	subject store.ID,
	emotion emotion.Emotion,
	message string,
) (*gqlmod.Reaction, *api.ResponseError) {
	t := h.c.t

	var result struct {
		CreateReaction *gqlmod.Reaction `json:"createReaction"`
	}
	err := h.c.QueryVar(
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
		map[string]string{
			"author":  string(author),
			"subject": string(subject),
			"emotion": string(emotion),
			"message": message,
		},
		&result,
	)

	if successAssumption {
		require.Nil(t, err, 0)
	} else if err != nil {
		return nil, err
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

	return result.CreateReaction, nil
}

// CreateReaction helps creating a reaction
func (h Helper) CreateReaction(
	author store.ID,
	subject store.ID,
	emotion emotion.Emotion,
	message string,
) (*gqlmod.Reaction, *api.ResponseError) {
	return h.createReaction(
		potentialFailure,
		author,
		subject,
		emotion,
		message,
	)
}

// CreateReaction helps creating a reaction and assumes success
func (ok AssumeSuccess) CreateReaction(
	author store.ID,
	subject store.ID,
	emotion emotion.Emotion,
	message string,
) *gqlmod.Reaction {
	result, _ := ok.h.createReaction(
		success,
		author,
		subject,
		emotion,
		message,
	)
	return result
}
