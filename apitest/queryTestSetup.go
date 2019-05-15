package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
)

type id = store.ID
type user = *gqlmod.User
type post = *gqlmod.Post
type reaction = *gqlmod.Reaction

type queryTestSetup struct {
	ts                *setup.TestSetup
	users             map[id]user
	posts             map[id]post
	reactions         map[id]reaction
	postsByAuthor     map[id]map[id]post
	authorByPosts     map[id]user
	reactionsByAuthor map[id]map[id]reaction
	authorByReactions map[id]user
	reactionsByPost   map[id]map[id]reaction
	subReaction       map[id]map[id]reaction
}

func (ts *queryTestSetup) Teardown() {
	ts.ts.Teardown()
}

func newQueryTestSetup(
	t *testing.T,
	tcx setup.TestContext,
) (s queryTestSetup) {
	ts := setup.New(t, tcx)
	s.ts = ts
	s.users = make(map[id]user)
	s.posts = make(map[id]post)
	s.reactions = make(map[id]reaction)
	s.postsByAuthor = make(map[id]map[id]post)
	s.authorByPosts = make(map[id]user)
	s.reactionsByAuthor = make(map[id]map[id]reaction)
	s.authorByReactions = make(map[id]user)
	s.reactionsByPost = make(map[id]map[id]reaction)
	s.subReaction = make(map[id]map[id]reaction)

	debug := ts.Debug()

	// Entities: users

	// User: first
	userA := debug.Help.OK.CreateUser(
		"first",
		"1@test.test",
		"testpass",
	)
	s.users[*userA.ID] = userA

	// User: second
	userB := debug.Help.OK.CreateUser(
		"second",
		"2@test.test",
		"testpass",
	)
	s.users[*userB.ID] = userB

	// User: third
	userC := debug.Help.OK.CreateUser(
		"third",
		"3@test.test",
		"testpass",
	)
	s.users[*userC.ID] = userC

	// Entities: posts

	// Post: A1
	postA1 := debug.Help.OK.CreatePost(
		*userA.ID,
		"Post A1",
		"Post A1 contents",
	)
	s.posts[*postA1.ID] = postA1

	// Post: A2
	postA2 := debug.Help.OK.CreatePost(
		*userA.ID,
		"Post A2",
		"Post A2 contents",
	)
	s.posts[*postA2.ID] = postA2

	// Post: B1
	postB1 := debug.Help.OK.CreatePost(
		*userB.ID,
		"Post B1",
		"Post B1 contents",
	)
	s.posts[*postB1.ID] = postB1

	// Entities: reactions

	// UserB --ReactionB1-> PostA1
	reactionB1 := debug.Help.OK.CreateReaction(
		*userB.ID,
		*postA1.ID,
		emotion.Happy,
		"Reaction B1",
	)
	s.reactions[*reactionB1.ID] = reactionB1

	// UserB --ReactionB2-> PostA2
	reactionB2 := debug.Help.OK.CreateReaction(
		*userB.ID,
		*postA2.ID,
		emotion.Excited,
		"Reaction B2",
	)
	s.reactions[*reactionB2.ID] = reactionB2

	// UserC --ReactionC1-> PostA1
	reactionC1 := debug.Help.OK.CreateReaction(
		*userC.ID,
		*postA1.ID,
		emotion.Thoughtful,
		"Reaction C1",
	)
	s.reactions[*reactionC1.ID] = reactionC1

	// UserC --ReactionC2-> ReactionB2
	reactionC2 := debug.Help.OK.CreateReaction(
		*userC.ID,
		*reactionB2.ID,
		emotion.Thoughtful,
		"Reaction C2",
	)
	s.reactions[*reactionC2.ID] = reactionC2

	// Index: author -> posts

	// User A -> posts
	s.postsByAuthor[*userA.ID] = make(map[id]post, 2)
	s.postsByAuthor[*userA.ID][*postA1.ID] = postA1
	s.postsByAuthor[*userA.ID][*postA2.ID] = postA2

	s.postsByAuthor[*userB.ID] = make(map[id]post, 1)
	s.postsByAuthor[*userB.ID][*postB1.ID] = postB1

	s.postsByAuthor[*userC.ID] = make(map[id]post)

	// Index: posts -> author
	for authorID, posts := range s.postsByAuthor {
		author := s.users[authorID]
		for postID := range posts {
			s.authorByPosts[postID] = author
		}
	}

	// Index: author -> reactions
	s.reactionsByAuthor[*userA.ID] = make(map[id]reaction)

	s.reactionsByAuthor[*userB.ID] = make(map[id]reaction)
	s.reactionsByAuthor[*userB.ID][*reactionB1.ID] = reactionB1
	s.reactionsByAuthor[*userB.ID][*reactionB2.ID] = reactionB2

	s.reactionsByAuthor[*userC.ID] = make(map[id]reaction)
	s.reactionsByAuthor[*userC.ID][*reactionC1.ID] = reactionC1
	s.reactionsByAuthor[*userC.ID][*reactionC2.ID] = reactionC2

	// Index: reactions -> author
	for authorID, reactions := range s.reactionsByAuthor {
		for reactionID := range reactions {
			s.authorByReactions[reactionID] = s.users[authorID]
		}
	}

	// Index: post -> reactions
	s.reactionsByPost[*postA1.ID] = make(map[id]reaction)
	s.reactionsByPost[*postA1.ID][*reactionB1.ID] = reactionB1
	s.reactionsByPost[*postA1.ID][*reactionC1.ID] = reactionC1

	s.reactionsByPost[*postA2.ID] = make(map[id]reaction)
	s.reactionsByPost[*postA2.ID][*reactionB2.ID] = reactionB2

	s.reactionsByPost[*postB1.ID] = make(map[id]reaction)

	// Index reaction -> reactions
	s.subReaction[*reactionB2.ID] = make(map[id]reaction)
	s.subReaction[*reactionB2.ID][*reactionC2.ID] = reactionC2
	return
}
