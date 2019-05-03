package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/stretchr/testify/require"
)

// TestQuery tests post creation
func TestQuery(t *testing.T) {
	t.Run("users", func(t *testing.T) {
		s := newQueryTestSetup(t, tcx)
		defer s.Teardown()

		clt := s.ts.Root()

		var query struct {
			Users []gqlmod.User `json:"users"`
		}
		clt.Query(
			`query {
				users {
					id
					creation
					displayName
					email
				}
			}`,
			&query,
		)
		require.Len(t, query.Users, len(s.users))
		for _, actual := range query.Users {
			require.Contains(t, s.users, *actual.ID)
			compareUsers(t, s.users[*actual.ID], &actual)
		}
	})

	t.Run("posts", func(t *testing.T) {
		s := newQueryTestSetup(t, tcx)
		defer s.Teardown()

		clt := s.ts.Root()

		var query struct {
			Posts []gqlmod.Post `json:"posts"`
		}
		clt.Query(
			`query {
				posts {
					id
					creation
					title
					contents
				}
			}`,
			&query,
		)
		require.Len(t, query.Posts, len(s.posts))
		for _, actual := range query.Posts {
			require.Contains(t, s.posts, *actual.ID)
			comparePosts(t, s.posts[*actual.ID], &actual)
		}
	})

	t.Run("user", func(t *testing.T) {
		s := newQueryTestSetup(t, tcx)
		defer s.Teardown()

		clt := s.ts.Root()

		for _, expected := range s.users {
			var query struct {
				User *gqlmod.User `json:"user"`
			}
			clt.QueryVar(
				`query($userId: Identifier!) {
					user(id: $userId) {
						id
						creation
						displayName
						email
					}
				}`,
				map[string]string{
					"userId": string(*expected.ID),
				},
				&query,
			)
			require.NotNil(t, query.User)
			compareUsers(t, expected, query.User)
		}
	})

	t.Run("user (inexistent)", func(t *testing.T) {
		s := newQueryTestSetup(t, tcx)
		defer s.Teardown()

		clt := s.ts.Root()

		var query struct {
			User *gqlmod.User `json:"user"`
		}
		clt.QueryVar(
			`query($userId: Identifier!) {
				user(id: $userId) {
					id
					creation
					displayName
					email
				}
			}`,
			map[string]string{
				"userId": string(store.NewID()),
			},
			&query,
		)
		require.Nil(t, query.User)
	})

	t.Run("post", func(t *testing.T) {
		s := newQueryTestSetup(t, tcx)
		defer s.Teardown()

		clt := s.ts.Root()

		for _, expected := range s.posts {
			var query struct {
				Post *gqlmod.Post `json:"post"`
			}
			clt.QueryVar(
				`query($postId: Identifier!) {
					post(id: $postId) {
						id
						creation
						title
						contents
					}
				}`,
				map[string]string{
					"postId": string(*expected.ID),
				},
				&query,
			)
			require.NotNil(t, query.Post)
			comparePosts(t, expected, query.Post)
		}
	})

	t.Run("User.posts", func(t *testing.T) {
		s := newQueryTestSetup(t, tcx)
		defer s.Teardown()

		clt := s.ts.Root()

		for authorID, posts := range s.postsByAuthor {
			var query struct {
				User *gqlmod.User `json:"user"`
			}
			clt.QueryVar(
				`query($userId: Identifier!) {
					user(id: $userId) {
						posts {
							id
							title
							contents
							creation
						}
					}
				}`,
				map[string]string{
					"userId": string(authorID),
				},
				&query,
			)

			require.NotNil(t, query.User)
			require.Len(t, query.User.Posts, len(posts))

			for _, actualPost := range query.User.Posts {
				id := *actualPost.ID
				require.Contains(t, posts, id)
				comparePosts(t, posts[id], &actualPost)
			}
		}
	})

	t.Run("Post.author", func(t *testing.T) {
		s := newQueryTestSetup(t, tcx)
		defer s.Teardown()

		clt := s.ts.Root()

		for postID, author := range s.authorByPosts {
			var query struct {
				Post *gqlmod.Post `json:"post"`
			}
			clt.QueryVar(
				`query($postId: Identifier!) {
					post(id: $postId) {
						author {
							id
							email
							displayName
							creation
						}
					}
				}`,
				map[string]string{
					"postId": string(postID),
				},
				&query,
			)

			require.NotNil(t, query.Post)
			require.NotNil(t, query.Post.Author)
			compareUsers(t, s.users[*author.ID], query.Post.Author)
		}
	})

	t.Run("User.sessions", func(t *testing.T) {
		s := newQueryTestSetup(t, tcx)
		defer s.Teardown()

		clt, sess := s.ts.Client("1@test.test", "testpass")

		type expectedSessions = []*gqlmod.Session

		expect := func(expectedSessions expectedSessions) {
			var query struct {
				User *gqlmod.User `json:"user"`
			}
			clt.QueryVar(
				`query($userId: Identifier!) {
					user(id: $userId) {
						sessions {
							key
							creation
						}
					}
				}`,
				map[string]string{
					"userId": string(*sess.User.ID),
				},
				&query,
			)

			require.NotNil(t, query.User)
			require.Len(t, query.User.Sessions, len(expectedSessions))

			// Create key -> session index
			index := make(map[string]*gqlmod.Session)
			for _, expected := range expectedSessions {
				index[*expected.Key] = expected
			}

			// Check actual sessions
			for _, actualSession := range query.User.Sessions {
				actualKey := *actualSession.Key
				require.Contains(t, index, actualKey)
			}
		}

		// Sign in twice
		expect(expectedSessions{sess})
		_, sess2 := s.ts.Client("1@test.test", "testpass")
		expect(expectedSessions{sess, sess2})
	})

	t.Run("Post.reactions", func(t *testing.T) {
		s := newQueryTestSetup(t, tcx)
		defer s.Teardown()

		clt := s.ts.Root()

		for postID, reactions := range s.reactionsByPost {
			var query struct {
				Post *gqlmod.Post `json:"post"`
			}
			clt.QueryVar(
				`query($postId: Identifier!) {
					post(id: $postId) {
						reactions {
							id
							emotion
							message
							creation
						}
					}
				}`,
				map[string]string{
					"postId": string(postID),
				},
				&query,
			)

			require.NotNil(t, query.Post)
			require.Len(t, query.Post.Reactions, len(reactions))

			for _, actualReaction := range query.Post.Reactions {
				id := *actualReaction.ID
				require.Contains(t, reactions, id)
				compareReactions(t, reactions[id], &actualReaction)
			}
		}
	})

	t.Run("User.publishedReactions", func(t *testing.T) {
		s := newQueryTestSetup(t, tcx)
		defer s.Teardown()

		clt := s.ts.Root()

		for authorID, reactions := range s.reactionsByAuthor {
			var query struct {
				User *gqlmod.User `json:"user"`
			}
			clt.QueryVar(
				`query($authorId: Identifier!) {
					user(id: $authorId) {
						publishedReactions {
							id
							emotion
							message
							creation
						}
					}
				}`,
				map[string]string{
					"authorId": string(authorID),
				},
				&query,
			)

			require.NotNil(t, query.User)
			require.Len(t, query.User.PublishedReactions, len(reactions))

			for _, actualReaction := range query.User.PublishedReactions {
				id := *actualReaction.ID
				require.Contains(t, reactions, id)
				compareReactions(t, reactions[id], &actualReaction)
			}
		}
	})

	t.Run("Reaction.reactions", func(t *testing.T) {
		s := newQueryTestSetup(t, tcx)
		defer s.Teardown()

		clt := s.ts.Root()

		for subjectReactionID, subReactions := range s.subReaction {
			var query struct {
				Reaction *gqlmod.Reaction `json:"reaction"`
			}
			clt.QueryVar(
				`query($subjectReactionId: Identifier!) {
					reaction(id: $subjectReactionId) {
						reactions {
							id
							emotion
							message
							creation
						}
					}
				}`,
				map[string]string{
					"subjectReactionId": string(subjectReactionID),
				},
				&query,
			)

			require.NotNil(t, query.Reaction)
			require.Len(t, query.Reaction.Reactions, len(subReactions))

			for _, actualReaction := range query.Reaction.Reactions {
				id := *actualReaction.ID
				require.Contains(t, subReactions, id)
				compareReactions(t, subReactions[id], &actualReaction)
			}
		}
	})
}
