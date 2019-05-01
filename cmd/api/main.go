package main

import (
	"context"
	"flag"
	"log"

	"github.com/romshark/dgraph_graphql_go/api"
	"github.com/romshark/dgraph_graphql_go/store"
)

var host = flag.String("host", "localhost:16000", "API server host address")
var dbHost = flag.String("dbhost", "localhost:9080", "database host address")

func main() {
	flag.Parse()

	str := store.NewStore(*dbHost)
	if err := str.Prepare(); err != nil {
		log.Fatalf("store prepare: %s", err)
	}

	api := api.NewServer(api.ServerOptions{
		Host: *host,
	}, str)

	if err := api.Launch(); err != nil {
		log.Fatalf("API server launch: %s", err)
	}

	// Setup termination signal listener
	onTerminate(func() {
		if err := api.Shutdown(context.Background()); err != nil {
			log.Fatalf("API server shutdown: %s", err)
		}
	})

	api.AwaitShutdown()
}
