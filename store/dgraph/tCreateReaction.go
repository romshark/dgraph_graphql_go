package dgraph

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/store"
	emo "github.com/romshark/dgraph_graphql_go/store/enum/emotion"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CreateReaction creates a new post
func (str *impl) CreateReaction(
	ctx context.Context,
	authorID store.ID,
	subjectID store.ID,
	emotion emo.Emotion,
	message string,
) (
	result struct {
		UID          string
		ID           store.ID
		SubjectUID   string
		AuthorUID    string
		CreationTime time.Time
	},
	err error,
) {
	// Validate input
	err = store.ValidateReactionMessage(message)
	if err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		return
	}
	err = emo.Validate(emotion)
	if err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		return
	}

	// Prepare
	result.ID = store.NewID()
	result.CreationTime = time.Now()

	// Begin transaction
	txn, close := str.txn(&err)
	if err != nil {
		return
	}
	defer close()

	// Ensure author and subject exist
	var res struct {
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
		&res,
	)
	if err != nil {
		return
	}

	if len(res.ByID) > 0 {
		err = errors.Errorf("duplicate Reaction.id: %s", result.ID)
		return
	}
	if len(res.Author) < 1 {
		err = strerr.Newf(
			strerr.ErrInvalidInput,
			"author not found",
		)
		return
	}
	// subjectType: "p" for post, "r" for reaction
	subjectType := "p"
	if len(res.PostSubject) > 0 {
		result.SubjectUID = res.PostSubject[0].NodeID
	} else if len(res.ReactionSubject) > 0 {
		subjectType = "r"
		result.SubjectUID = res.ReactionSubject[0].NodeID
	} else {
		err = strerr.Newf(
			strerr.ErrInvalidInput,
			"subject not found",
		)
		return
	}

	result.AuthorUID = res.Author[0].NodeID

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
		Author:   UID{NodeID: result.AuthorUID},
		Subject:  UID{NodeID: result.SubjectUID},
		Emotion:  string(emotion),
		Message:  message,
		Creation: result.CreationTime,
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
		UID:                result.AuthorUID,
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
			UID:       result.SubjectUID,
			Reactions: UID{NodeID: result.UID},
		})
	} else {
		// Update reaction (Reaction.reactions -> new reaction)
		updateSubjectJSON, err = json.Marshal(struct {
			UID       string `json:"uid"`
			Reactions UID    `json:"Reaction.reactions"`
		}{
			UID:       result.SubjectUID,
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
