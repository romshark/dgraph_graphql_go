package http

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

// handleDebugAuth handles a debug authentication request
func (t *Server) handleDebugAuth(
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

	debugSessionKey := t.onDebugSess(req.Context(), pair[0], pair[1])
	if debugSessionKey == nil {
		unauthorized()
		return
	}

	// Set the session key as cookie
	http.SetCookie(resp, &http.Cookie{
		Name:     "SID",
		Value:    string(debugSessionKey),
		HttpOnly: true,
		Secure:   true,
	})

	// Return session key
	if _, err := resp.Write(debugSessionKey); err != nil {
		log.Printf("writing debug auth response: %s", err)
	}
}
