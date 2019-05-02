package api

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

// handleRoot handles an HTTP root signin request
func (srv *server) handleRoot(
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

	if srv.opts.RootUser.Status == RootUserDisabled {
		unauthorized()
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

	// Check root credentials
	if pair[0] != srv.opts.RootUser.Username ||
		pair[1] != srv.opts.RootUser.Password {
		unauthorized()
		return
	}

	// Set the session key as cookie
	http.SetCookie(resp, &http.Cookie{
		Name:     "SID",
		Value:    string(srv.rootSessionKey),
		HttpOnly: true,
		Secure:   true,
	})

	// Return session key
	if _, err := resp.Write(srv.rootSessionKey); err != nil {
		log.Printf("writing /root response: %s", err)
	}
}
