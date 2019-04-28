package main

import (
	"context"
	"flag"
	"log"
	"time"

	"demo/apisrv"
	"demo/store"
	"demo/store/dbmod"
)

var host = flag.String("host", "localhost:16000", "API server host address")
var dbHost = flag.String("dbhost", "localhost:9080", "database host address")

func main() {
	flag.Parse()

	str := store.NewStore(*dbHost)
	if err := str.Prepare(); err != nil {
		log.Fatalf("store prepare: %s", err)
	}

	apiSrv := apisrv.NewAPIServer(apisrv.APIServerOptions{
		Host: *host,
	}, str)

	if err := apiSrv.Launch(); err != nil {
		log.Fatalf("API server launch: %s", err)
	}

	// Setup termination signal listener
	onTerminate(func() {
		if err := apiSrv.Shutdown(context.Background()); err != nil {
			log.Fatalf("API server shutdown: %s", err)
		}
	})

	time.AfterFunc(time.Millisecond*50, func() {
		ctx := context.Background()
		aliceID, err := str.CreateUser(ctx, "alice@robinsons.net", "alice")
		if err != nil {
			log.Fatalf("alice creation: %s", err)
		}

		var aliceRes struct {
			Alice []dbmod.User `json:"alice"`
		}
		str.QueryVars(ctx, `
			query Alice($id: string) {
				alice(func: eq(User.id, $id)) {
					uid
					User.displayName
					User.email
				}
			}
		`, map[string]string{
			"$id": string(aliceID),
		}, &aliceRes)
	})

	apiSrv.AwaitShutdown()
}
