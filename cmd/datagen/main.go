package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/romshark/dgraph_graphql_go/datagen"
	gqlgen "github.com/romshark/dgraph_graphql_go/datagen/graphql"
)

var apiHost = flag.String("apihost", "localhost:16000", "API host address")
var debugUsername = flag.String("debug-username", "debug", "debug username")
var debugPassword = flag.String("debug-password", "debug", "debug password")
var verbose = flag.Bool("verbose", true, "verbose mode")
var concurrency = flag.Uint64("par", 1, "maximum concurrent mutations")

func main() {
	flag.Parse()

	// Initialize generator
	statisticsRecorder := datagen.NewStatisticsRecorder()
	generator, err := gqlgen.New(gqlgen.Config{
		APIHost: url.URL{
			Scheme: "https",
			Host:   *apiHost,
			Path:   "/g",
		},
		DebugUsername: *debugUsername,
		DebugPassword: *debugPassword,
		Verbose:       *verbose,
		Concurrency:   uint32(*concurrency),
	})
	if err != nil {
		log.Fatalf("graphql datagen init: %s", err)
	}

	// Load options
	options, err := datagen.OptionsFromFile("options.json")
	if err != nil {
		log.Fatalf("options loading: %s", err)
	}

	// Generate
	if _, err := generator.Generate(
		context.Background(),
		options,
		statisticsRecorder,
	); err != nil {
		log.Fatalf("generator: %s", err)
	}

	// Print statistics
	stats := datagen.StatisticsReader(statisticsRecorder)
	log.Print("generation successful:")
	fmt.Printf(
		" users: %d (%s; avg: %s)\n",
		stats.Users(),
		stats.UsersCreation(),
		stats.UserCreationAvg(),
	)
	fmt.Printf(
		" posts: %d (%s; avg: %s)\n",
		stats.Posts(),
		stats.PostsCreation(),
		stats.PostCreationAvg(),
	)
}
