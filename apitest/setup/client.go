package setup

import (
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/stretchr/testify/require"
)

// Client represents an API client
type Client struct {
	t              *testing.T
	ts             *TestSetup
	httpClient     *http.Client
	rootSessionKey []byte

	Help Helper
}

// Guest creates a new unauthenticated API client
func (ts *TestSetup) Guest() *Client {
	// Initialize the http cookie jar
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		ts.t.Fatalf("client cookie jar init: %s", err)
	}

	// Initialize client
	clt := &Client{
		t:  ts.t,
		ts: ts,
		httpClient: &http.Client{
			Timeout: time.Second * 1,
			Jar:     cookieJar,
		},
	}

	// Initialize helper
	clt.Help = Helper{
		c:                      clt,
		creationTimeTollerance: time.Second * 3,
	}
	clt.Help.OK = AssumeSuccess{
		h: &clt.Help,
		t: ts.t,
	}

	return clt
}

// Client creates a new authenticated API client
func (ts *TestSetup) Client(
	email,
	password string,
) (*Client, *gqlmod.Session) {
	clt := ts.Guest()

	sess, err := clt.Help.SignIn(email, password)
	require.Nil(ts.t, err)

	return clt, sess
}
