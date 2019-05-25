package graphql

import (
	"net/url"

	"github.com/romshark/dgraph_graphql_go/api/transport"
	thttp "github.com/romshark/dgraph_graphql_go/api/transport/http"
	"github.com/romshark/dgraph_graphql_go/datagen"
	"golang.org/x/sync/semaphore"
)

// gqlgen implements a random data generator for API testing purposes
type gqlgen struct {
	config          Config
	client          transport.Client
	slots           *semaphore.Weighted
	emailGen        *emailAddressGenerator
	displayNameGen  *displayNameGenerator
	postTitleGen    *postTitleGenerator
	postContentsGen *postContentsGenerator
}

// Config defines the data generator configuration
type Config struct {
	APIHost       url.URL
	Verbose       bool
	DebugUsername string
	DebugPassword string
	Concurrency   uint32
}

// New creates a new GraphQL-API-based data generator instance
func New(config Config) (datagen.DataGenerator, error) {
	// Initialize API client
	client, err := thttp.NewClient(config.APIHost, thttp.ClientOptions{})
	if err != nil {
		return nil, err
	}

	// Sign in API client as debug user
	if err := client.SignInDebug(
		config.DebugUsername,
		config.DebugPassword,
	); err != nil {
		return nil, err
	}

	return &gqlgen{
		config:          config,
		client:          client,
		slots:           semaphore.NewWeighted(int64(config.Concurrency)),
		emailGen:        newEmailAddressGenerator(),
		displayNameGen:  newDisplayNameGenerator(),
		postTitleGen:    newPostTitleGenerator(),
		postContentsGen: newPostContentsGenerator(256),
	}, nil
}
