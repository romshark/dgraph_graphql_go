package http

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

// handleRootAuth handles a root authentication request
func (t *Server) handleRootAuth(
	resp http.ResponseWriter,
	req *http.Request,
) {
	unauthorized := func() {
		http.Error(
			resp,
			http.StatusText(http.StatusForbidden),
			http.StatusForbidden,
		)
	}

	// Parse credentials
	s := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		unauthorized()
		return
	}
	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		unauthorized()
		return
	}
	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		unauthorized()
		return
	}

	rootSessionKey := t.onRootSess(req.Context(), pair[0], pair[1])
	if rootSessionKey == nil {
		unauthorized()
		return
	}

	// Set the session key as cookie
	http.SetCookie(resp, &http.Cookie{
		Name:     "SID",
		Value:    string(rootSessionKey),
		HttpOnly: true,
		Secure:   true,
	})

	// Return session key
	if _, err := resp.Write(rootSessionKey); err != nil {
		log.Printf("writing root auth response: %s", err)
	}
}
