package setup

import (
	"net/url"
	"testing"
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph"
	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	trn "github.com/romshark/dgraph_graphql_go/api/transport"
	thttp "github.com/romshark/dgraph_graphql_go/api/transport/http"
	"github.com/stretchr/testify/require"
)

// Client represents an API client
type Client struct {
	t          *testing.T
	ts         *TestSetup
	apiClient  trn.Client
	sessionKey []byte

	Help Helper
}

// Query performs an API query
func (tclt *Client) Query(
	query string,
	result interface{},
) error {
	return tclt.apiClient.Query(query, result)
}

// QueryVar performs a parameterized API query
func (tclt *Client) QueryVar(
	query string,
	vars map[string]interface{},
	result interface{},
) error {
	return tclt.apiClient.QueryVar(query, vars, result)
}

// Guest creates a new unauthenticated API client
func (ts *TestSetup) Guest() *Client {
	// Initialize client
	apiClt, err := thttp.NewClient(
		url.URL{
			Scheme: "http",
			Host:   ts.serverTransport.(*thttp.Server).Addr().Host,
		},
		thttp.ClientOptions{
			Timeout: time.Second * 10,
		},
	)
	require.NoError(ts.t, err)

	clt := &Client{
		t:         ts.t,
		ts:        ts,
		apiClient: apiClt,
	}

	// Initialize helper
	clt.Help = Helper{
		ts:                     ts,
		c:                      clt,
		creationTimeTollerance: time.Second * 3,
	}
	clt.Help.OK = AssumeSuccess{
		h: &clt.Help,
		t: ts.t,
	}

	return clt
}

// Debug creates a new authenticated API debug client
func (ts *TestSetup) Debug() *Client {
	clt := ts.Guest()

	// Sign in as debug user
	require.NoError(ts.t, clt.apiClient.AuthDebug(
		ts.debugUsername,
		ts.debugPassword,
	))

	return clt
}

// Client creates a new authenticated API client
func (ts *TestSetup) Client(
	email,
	password string,
) (*Client, *gqlmod.Session) {
	clt := ts.Guest()

	sess, err := clt.apiClient.Auth(email, password)
	require.Nil(ts.t, err)

	return clt, sess
}

func checkErr(
	t *testing.T,
	assumedSuccess successAssumption,
	err error,
) *graph.ResponseError {
	if !assumedSuccess {
		if err != nil {
			// In case of expected errors the error must be a graph error
			require.IsType(t, &graph.ResponseError{}, err)
			return err.(*graph.ResponseError)
		}
		return nil
	}
	require.NoError(t, err)
	return nil
}
