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
	isDebug    bool
	sessionKey string
}

// NewClient creates a new API client instance
func NewClient(host url.URL, conf ClientConfig) (trn.Client, error) {
	// Initialize the http cookie jar
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "client cookie jar init")
	}

	// Initialize client
	return &Client{
		host: host,
		httpClt: &http.Client{
			Timeout: conf.Timeout,
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
		if c.isDebug {
			req.Header.Set("Authorization", "Debug "+c.sessionKey)
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
		Data  interface{}         `json:"data"`
		Error *graphResponseError `json:"errors"`
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

// SignIn implements the transport.Client interface
func (c *Client) SignIn(email, password string) (*gqlmod.Session, error) {
	var result struct {
		CreateSession gqlmod.Session `json:"createSession"`
	}
	if err := c.QueryVar(
		`mutation(
			$email: String!
			$password: String!
		) {
			createSession(
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

	c.sessionKey = *result.CreateSession.Key

	return &result.CreateSession, nil
}

// SignInDebug implements the transport.Client interface
func (c *Client) SignInDebug(username, password string) error {
	// Sign in as debug
	// Initialize request
	u := c.host
	u.Path = "/debug"
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return errors.Wrap(err, "POST /debug request creation")
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
		return errors.Wrap(err, "POST /debug request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf(
			"debug signin bad response code: %d",
			resp.StatusCode,
		)
	}

	debugSessionKey, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read POST /debug response")
	}
	c.isDebug = true
	c.sessionKey = string(debugSessionKey)

	return nil
}

// Auth implements the transport.Client interface
func (c *Client) Auth(sessionKey string) (*gqlmod.Session, error) {
	var result struct {
		Authenticate gqlmod.Session `json:"authenticate"`
	}
	if err := c.QueryVar(
		`mutation(
			$sessionKey: String!
		) {
			authenticate(sessionKey: $sessionKey) {
				key
				creation
				user {
					id
				}
			}
		}`,
		map[string]interface{}{
			"sessionKey": sessionKey,
		},
		&result,
	); err != nil {
		return nil, err
	}

	c.sessionKey = *result.Authenticate.Key

	return &result.Authenticate, nil
}
