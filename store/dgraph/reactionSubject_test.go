package dgraph_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/dgraph"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
	"github.com/stretchr/testify/require"
)

// TestReactionSubjectUnMarshal tests marshalling and unmarshalling of
// ReactionSubject unions
func TestReactionSubjectUnMarshal(t *testing.T) {
	timeNow := time.Now()
	timeNow = timeNow.Truncate(time.Second)

	// Post <-> JSON
	t.Run("JSON_Post", func(t *testing.T) {
		post := dgraph.Post{
			UID:      "0x01",
			ID:       store.NewID(),
			Creation: timeNow,
			Title:    "Test title",
			Contents: "Test contents",
		}
		marshedPost, err := json.Marshal(dgraph.ReactionSubject{&post})
		require.NoError(t, err)

		var u dgraph.ReactionSubject
		require.NoError(t, json.Unmarshal(marshedPost, &u))
		require.IsType(t, &post, u.V)
		require.Equal(t, post, *u.V.(*dgraph.Post))
	})

	// Reaction <-> JSON
	t.Run("JSON_Reaction", func(t *testing.T) {
		reaction := dgraph.Reaction{
			UID:      "0x01",
			ID:       store.NewID(),
			Creation: timeNow,
			Emotion:  emotion.Excited,
			Message:  "Test message",
		}
		marshedReaction, err := json.Marshal(dgraph.ReactionSubject{&reaction})
		require.NoError(t, err)

		var u dgraph.ReactionSubject
		require.NoError(t, json.Unmarshal(marshedReaction, &u))
		require.IsType(t, &reaction, u.V)
		require.Equal(t, reaction, *u.V.(*dgraph.Reaction))
	})

	// Invalid type <-> JSON
	t.Run("JSON_invalid", func(t *testing.T) {
		require.Panics(t, func() {
			json.Marshal(dgraph.ReactionSubject{"invalid_union_type"})
		})
		invalidJson := []byte(`{"Foo": "bar", "Baz": 42}`)
		var u dgraph.ReactionSubject
		require.Error(t, json.Unmarshal(invalidJson, &u))
	})
}
