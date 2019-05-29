package main

import (
	"context"
	"flag"
	"log"

	"github.com/romshark/dgraph_graphql_go/api"
	"github.com/romshark/dgraph_graphql_go/api/config"
)

var argConfigFile = flag.String(
	"config",
	"./config.toml",
	"path to the configuration file",
)

func main() {
	flag.Parse()

	serverConfig, err := config.FromFile(*argConfigFile)
	if err != nil {
		log.Fatalf("reading config: %s", err)
	}

	api, err := api.NewServer(serverConfig)
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
