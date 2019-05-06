package main

import (
	"context"
	"flag"
	"log"

	"github.com/romshark/dgraph_graphql_go/api"
	"github.com/romshark/dgraph_graphql_go/api/passhash"
	"github.com/romshark/dgraph_graphql_go/api/sesskeygen"
	"github.com/romshark/dgraph_graphql_go/api/transport"
	thttp "github.com/romshark/dgraph_graphql_go/api/transport/http"
)

var host = flag.String("host", "localhost:16000", "API server host address")
var dbHost = flag.String("dbhost", "localhost:9080", "database host address")

func main() {
	flag.Parse()

	api, err := api.NewServer(api.ServerOptions{
		Host:                *host,
		DBHost:              *dbHost,                 // database host address
		SessionKeyGenerator: sesskeygen.NewDefault(), // session key generator
		PasswordHasher:      passhash.Bcrypt{},       // password hasher
		Transport: []transport.Server{
			thttp.NewServer(thttp.ServerOptions{
				Host: *host,
			}), // HTTP transport
		},
	})
	if err != nil {
		log.Fatalf("API server init: %s", err)
	}

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
