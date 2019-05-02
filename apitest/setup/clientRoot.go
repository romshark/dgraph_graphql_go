package setup

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

// Root creates a new authenticated API root client
func (ts *TestSetup) Root() *Client {
	clt := ts.Guest()

	// Sign in as root
	// Initialize request
	u := url.URL{
		Scheme: "http",
		Host:   ts.apiServer.Addr().Host,
		Path:   "/root",
	}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		ts.t.Fatalf("POST /root request creation: %s", err)
	}

	// Set authentication header
	req.Header.Add("Authorization", ts.rootHTTPBasicAuth)

	// Perform request
	resp, err := clt.httpClient.Do(req)
	if err != nil {
		ts.t.Fatalf("POST /root request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ts.t.Fatalf("root signin bad response code: %d", resp.StatusCode)
	}

	rootSessionKey, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ts.t.Fatalf("read POST /root response: %s", err)
	}
	clt.rootSessionKey = rootSessionKey

	return clt
}
