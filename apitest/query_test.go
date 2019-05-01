package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/stretchr/testify/require"
)

// TestQuery tests post creation
func TestQuery(t *testing.T) {
	type TestSetup struct {
		ts            *setup.TestSetup
		users         map[store.ID]*gqlmod.User
		posts         map[store.ID]*gqlmod.Post
		postsByAuthor map[store.ID]map[store.ID]*gqlmod.Post
	}

	setupTest := func(t *testing.T) TestSetup {
		ts := setup.New(t, tcx)

		userA := ts.Help.OK.CreateUser("fooBarowich", "foo@bar.buz")
		userB := ts.Help.OK.CreateUser("buzBazowich", "buz@foo.foo")
		userC := ts.Help.OK.CreateUser("fuzFuzzowich", "fuz@fuz.fuz")

		postA1 := ts.Help.OK.CreatePost(*userA.ID, "A post 1", "test content 1")
		postA2 := ts.Help.OK.CreatePost(*userA.ID, "A post 2", "test content 2")
		postB1 := ts.Help.OK.CreatePost(*userB.ID, "B post 1", "test content 3")

		users := make(map[store.ID]*gqlmod.User, 2)
		users[*userA.ID] = userA
		users[*userB.ID] = userB

		posts := make(map[store.ID]*gqlmod.Post, 3)
		posts[*postA1.ID] = postA1
		posts[*postA2.ID] = postA2
		posts[*postB1.ID] = postB1

		// Index: posts by author
		postsByAuthor := make(
			map[store.ID]map[store.ID]*gqlmod.Post,
			len(posts),
		)

		postsByAuthor[*userA.ID] = make(map[store.ID]*gqlmod.Post, 2)
		postsByAuthor[*userA.ID][*postA1.ID] = postA1
		postsByAuthor[*userA.ID][*postA2.ID] = postA2

		postsByAuthor[*userB.ID] = make(map[store.ID]*gqlmod.Post, 1)
		postsByAuthor[*userB.ID][*postB1.ID] = postB1

		postsByAuthor[*userC.ID] = make(map[store.ID]*gqlmod.Post)

		return TestSetup{
			ts:            ts,
			users:         users,
			posts:         posts,
			postsByAuthor: postsByAuthor,
		}
	}

	t.Run("users", func(t *testing.T) {
		s := setupTest(t)
		defer s.ts.Teardown()

		var query struct {
			Users []gqlmod.User `json:"users"`
		}
		s.ts.Query(
			`query {
				users {
					id
					creation
					displayName
					email
					posts {
						id
					}
				}
			}`,
			&query,
		)
		require.Len(t, query.Users, len(s.users))
		for _, actual := range query.Users {
			require.Contains(t, s.users, *actual.ID)
			require.Equal(t, *s.users[*actual.ID], actual)
		}
	})

	t.Run("posts", func(t *testing.T) {
		s := setupTest(t)
		defer s.ts.Teardown()

		var query struct {
			Posts []gqlmod.Post `json:"posts"`
		}
		s.ts.Query(
			`query {
				posts {
					id
					creation
					title
					contents
					author {
						id
					}
				}
			}`,
			&query,
		)
		require.Len(t, query.Posts, len(s.posts))
		for _, actual := range query.Posts {
			require.Contains(t, s.posts, *actual.ID)
			require.Equal(t, *s.posts[*actual.ID], actual)
		}
	})

	t.Run("user", func(t *testing.T) {
		s := setupTest(t)
		defer s.ts.Teardown()

		for _, expected := range s.users {
			var query struct {
				User *gqlmod.User `json:"user"`
			}
			s.ts.QueryVar(
				`query($userId: Identifier!) {
					user(id: $userId) {
						id
						creation
						displayName
						email
						posts {
							id
						}
					}
				}`,
				map[string]string{
					"userId": string(*expected.ID),
				},
				&query,
			)
			require.NotNil(t, query.User)
			require.Equal(t, *expected, *query.User)
		}
	})

	t.Run("post", func(t *testing.T) {
		s := setupTest(t)
		defer s.ts.Teardown()

		for _, expected := range s.posts {
			var query struct {
				Post *gqlmod.Post `json:"post"`
			}
			s.ts.QueryVar(
				`query($postId: Identifier!) {
					post(id: $postId) {
						id
						creation
						title
						contents
						author {
							id
						}
					}
				}`,
				map[string]string{
					"postId": string(*expected.ID),
				},
				&query,
			)
			require.NotNil(t, query.Post)
			require.Equal(t, *expected, *query.Post)
		}
	})

	t.Run("user.posts", func(t *testing.T) {
		s := setupTest(t)
		defer s.ts.Teardown()

		for authorID, posts := range s.postsByAuthor {
			var query struct {
				User *gqlmod.User `json:"user"`
			}
			s.ts.QueryVar(
				`query($userId: Identifier!) {
					user(id: $userId) {
						posts {
							id
							title
							contents
							creation
							author {
								id
							}
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
				require.Equal(t, *posts[id], actualPost)
			}
		}
	})
}
