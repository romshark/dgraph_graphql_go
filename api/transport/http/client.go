package http

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/graph"
	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	trn "github.com/romshark/dgraph_graphql_go/api/transport"
)

// Client represents an HTTP client implementation
type Client struct {
	host       url.URL
	httpClt    *http.Client
	isRoot     bool
	sessionKey string
}

// NewClient creates a new API client instance
func NewClient(host url.URL, opts ClientOptions) (trn.Client, error) {
	// Initialize the http cookie jar
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "client cookie jar init")
	}

	// Initialize client
	return &Client{
		host: host,
		httpClt: &http.Client{
			Timeout: opts.Timeout,
			Jar:     cookieJar,
		},
	}, nil
}

// Query implements the transport.Client interface
func (c *Client) Query(
	query string,
	result interface{},
) error {
	return c.QueryVar(query, nil, result)
}

// QueryVar implements the transport.Client interface
func (c *Client) QueryVar(
	query string,
	vars map[string]interface{},
	result interface{},
) error {
	// Marshal form data
	requestData := struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}{
		Query:     query,
		Variables: vars,
	}
	marshed, err := json.Marshal(&requestData)
	if err != nil {
		return errors.Wrap(err, "query marshal")
	}

	u := c.host
	u.Path = "/g"

	// Initialize request
	req, err := http.NewRequest(
		"POST",
		u.String(),
		bytes.NewBuffer(marshed),
	)
	if err != nil {
		return errors.Wrap(err, "query POST request creation")
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Set authorization headers if authentication
	if c.sessionKey != "" {
		if c.isRoot {
			req.Header.Set("Authorization", "Root "+c.sessionKey)
		} else {
			req.Header.Set("Authorization", "Bearer "+c.sessionKey)
		}
	}

	// Perform request
	resp, err := c.httpClt.Do(req)
	if err != nil {
		return errors.Wrap(err, "HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusBadRequest {
		return errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	responseDecoderJSON := json.NewDecoder(resp.Body)

	res := struct {
		Data  interface{}         `json:"d"`
		Error *graphResponseError `json:"e"`
	}{
		Data: result,
	}
	if err := responseDecoderJSON.Decode(&res); err != nil {
		return errors.Wrap(err, "response decode JSON")
	}

	if res.Error != nil {
		return &graph.ResponseError{
			Code:    res.Error.Code,
			Message: res.Error.Message,
		}
	}

	return nil
}

// Auth implements the transport.Client interface
func (c *Client) Auth(email, password string) (*gqlmod.Session, error) {
	var result struct {
		SignIn gqlmod.Session `json:"signIn"`
	}
	if err := c.QueryVar(
		`mutation(
			$email: String!
			$password: String!
		) {
			signIn(
				email: $email
				password: $password
			) {
				key
				user {
					id
				}
				creation
			}
		}`,
		map[string]interface{}{
			"email":    email,
			"password": password,
		},
		&result,
	); err != nil {
		return nil, err
	}

	c.sessionKey = *result.SignIn.Key

	return &result.SignIn, nil
}

// AuthRoot implements the transport.Client interface
func (c *Client) AuthRoot(username, password string) error {
	// Sign in as root
	// Initialize request
	u := c.host
	u.Path = "/root"
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return errors.Wrap(err, "POST /root request creation")
	}

	// Set authentication header
	req.Header.Add(
		"Authorization",
		"Basic "+base64.StdEncoding.EncodeToString(
			[]byte(username+":"+password),
		),
	)

	// Perform request
	resp, err := c.httpClt.Do(req)
	if err != nil {
		return errors.Wrap(err, "POST /root request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf(
			"root signin bad response code: %d",
			resp.StatusCode,
		)
	}

	rootSessionKey, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read POST /root response")
	}
	c.isRoot = true
	c.sessionKey = string(rootSessionKey)

	return nil
}
