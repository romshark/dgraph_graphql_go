package dgraph

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/store"
	emo "github.com/romshark/dgraph_graphql_go/store/enum/emotion"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CreateReaction creates a new post
func (str *impl) CreateReaction(
	ctx context.Context,
	creationTime time.Time,
	authorID store.ID,
	subjectID store.ID,
	emotion emo.Emotion,
	message string,
) (
	result store.Reaction,
	err error,
) {
	result.Creation = creationTime
	result.Emotion = emotion
	result.Message = message

	// Prepare
	result.ID = store.NewID()

	// Begin transaction
	txn, close := str.txn(&err)
	if err != nil {
		return
	}
	defer close()

	// Ensure author and subject exist
	var qr struct {
		ByID            []UID `json:"byId"`
		Author          []UID `json:"author"`
		PostSubject     []UID `json:"postSubject"`
		ReactionSubject []UID `json:"reactionSubject"`
	}
	err = txn.QueryVars(
		ctx,
		`query User(
			$id: string,
			$authorId: string,
			$subjectId: string
		) {
			byId(func: eq(Reaction.id, $id)) { uid }
			author(func: eq(User.id, $authorId)) { uid }
			postSubject(func: eq(Post.id, $subjectId)) { uid }
			reactionSubject(func: eq(Reaction.id, $subjectId)) { uid }
		}`,
		map[string]string{
			"$id":        string(result.ID),
			"$authorId":  string(authorID),
			"$subjectId": string(subjectID),
		},
		&qr,
	)
	if err != nil {
		return
	}

	if len(qr.ByID) > 0 {
		err = errors.Errorf("duplicate Reaction.id: %s", result.ID)
		return
	}
	if len(qr.Author) < 1 {
		err = strerr.Newf(
			strerr.ErrInvalidInput,
			"author not found",
		)
		return
	}
	// subjectType: "p" for post, "r" for reaction
	subjectType := "p"
	if len(qr.PostSubject) > 0 {
		result.Subject = store.Post{
			GraphNode: store.GraphNode{
				UID: qr.PostSubject[0].NodeID,
			},
		}
	} else if len(qr.ReactionSubject) > 0 {
		subjectType = "r"
		result.Subject = store.Reaction{
			GraphNode: store.GraphNode{
				UID: qr.ReactionSubject[0].NodeID,
			},
		}
	} else {
		err = strerr.Newf(
			strerr.ErrInvalidInput,
			"subject not found",
		)
		return
	}

	result.Author = &store.User{
		GraphNode: store.GraphNode{
			UID: qr.Author[0].NodeID,
		},
	}

	// Create new reaction
	var newReactionJSON []byte
	newReactionJSON, err = json.Marshal(struct {
		ID       string    `json:"Reaction.id"`
		Author   UID       `json:"Reaction.author"`
		Subject  UID       `json:"Reaction.subject"`
		Emotion  string    `json:"Reaction.emotion"`
		Message  string    `json:"Reaction.message"`
		Creation time.Time `json:"Reaction.creation"`
	}{
		ID:       string(result.ID),
		Author:   UID{NodeID: result.Author.UID},
		Subject:  UID{NodeID: result.Subject.NodeID()},
		Emotion:  string(emotion),
		Message:  message,
		Creation: creationTime,
	})
	if err != nil {
		return
	}
	var reactionCreationMut map[string]string
	reactionCreationMut, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: newReactionJSON,
	})
	if err != nil {
		return
	}
	result.UID = reactionCreationMut["blank-0"]

	// Update author (User.publishedReactions -> new reaction)
	var updatedAuthorJSON []byte
	updatedAuthorJSON, err = json.Marshal(struct {
		UID                string `json:"uid"`
		PublishedReactions UID    `json:"User.publishedReactions"`
	}{
		UID:                result.Author.UID,
		PublishedReactions: UID{NodeID: result.UID},
	})
	if err != nil {
		return
	}

	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: updatedAuthorJSON,
	})
	if err != nil {
		return
	}

	// Update subject
	var updateSubjectJSON []byte
	if subjectType == "p" {
		// Update post (Post.reactions -> new reaction)
		updateSubjectJSON, err = json.Marshal(struct {
			UID       string `json:"uid"`
			Reactions UID    `json:"Post.reactions"`
		}{
			UID:       result.Subject.NodeID(),
			Reactions: UID{NodeID: result.UID},
		})
	} else {
		// Update reaction (Reaction.reactions -> new reaction)
		updateSubjectJSON, err = json.Marshal(struct {
			UID       string `json:"uid"`
			Reactions UID    `json:"Reaction.reactions"`
		}{
			UID:       result.Subject.NodeID(),
			Reactions: UID{NodeID: result.UID},
		})
	}
	if err != nil {
		return
	}
	_, err = txn.Mutation(ctx, &api.Mutation{
		SetJson: updateSubjectJSON,
	})
	if err != nil {
		return
	}

	return
}
