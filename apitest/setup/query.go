package setup

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/romshark/dgraph_graphql_go/api"
)

// Query performs a query on the test API
func (c *Client) Query(
	query string,
	result interface{},
) *api.ResponseError {
	return c.QueryVar(query, nil, result)
}

// QueryVar performs a parameterized query on the test API
func (c *Client) QueryVar(
	query string,
	vars map[string]string,
	result interface{},
) *api.ResponseError {
	// Marshal form data
	requestData := struct {
		Query         string            `json:"query"`
		OperationName string            `json:"operationName"`
		Variables     map[string]string `json:"variables"`
	}{
		Query:     query,
		Variables: vars,
	}
	marshed, err := json.Marshal(&requestData)
	if err != nil {
		c.ts.t.Fatalf("query marshal: %s", err)
	}

	// Initialize request
	u := url.URL{
		Scheme: "http",
		Host:   c.ts.apiServer.Addr().Host,
		Path:   "/g",
	}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(marshed))
	if err != nil {
		c.ts.t.Fatalf("query POST request creation: %s", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.t.Fatalf("HTTP request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusBadRequest {
		c.t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.t.Fatalf("read HTTP response: %s", err)
	}

	var res api.Response
	res.Data = result
	if err := json.Unmarshal(body, &res); err != nil {
		c.t.Fatalf("unmarshal response: %s", err)
	}

	return res.Error
}
