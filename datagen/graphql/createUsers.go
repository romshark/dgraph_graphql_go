package graphql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/datagen"
	"github.com/schollz/progressbar"
)

func (gen *gqlgen) createUsers(
	ctx context.Context,
	graph *datagen.Graph,
	stats datagen.StatisticsWriter,
) error {
	wg := NewErrWaitGroup(uint64(len(graph.Users)))

	// Create user profile
	fmt.Println("")
	log.Println("create users")
	progressUserCrt := progressbar.New(len(graph.Users))
	progressUserCrt.RenderBlank()
	for i, u := range graph.Users {
		index := i
		user := u
		go func() {
			gen.slots.Acquire(ctx, 1)
			defer gen.slots.Release(1)

			start := time.Now()
			var qr struct {
				NewUser gqlmod.User `json:"createUser"`
			}
			if err := gen.client.QueryVar(
				`mutation(
					$email: String!
					$displayName: String!
					$password: String!
				) {
					createUser(
						email: $email,
						displayName: $displayName,
						password: $password
					) {
						id
						email
						displayName
						creation
					}
				}`,
				map[string]interface{}{
					"email":       *user.Email,
					"displayName": *user.DisplayName,
					"password":    graph.UserPasswords[index],
				},
				&qr,
			); err != nil {
				wg.Fail(err)
				return
			}
			graph.Users[index] = qr.NewUser
			stats.RecordUserCreation(time.Since(start))
			progressUserCrt.Add(1)
			wg.Dec(1)
		}()
	}

	return wg.Wait()
}
