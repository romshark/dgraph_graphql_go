package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
)

// auth reads the session key from the authorization header, approaches the API
// server in order to verify the session key and if a session is returned it
// moves it to the request context
func (t *Server) auth(req *http.Request) *http.Request {
	// Set default (empty) session
	session := &auth.RequestSession{
		ShieldClientRole: auth.GQLShieldClientGuest,
	}
	req = req.WithContext(context.WithValue(
		req.Context(),
		auth.CtxSession,
		session,
	))

	// Try read the HTTP Authorization header
	authHeader := req.Header.Get("Authorization")
	if len(authHeader) < 1 {
		return req
	}

	tokens := strings.Split(authHeader, " ")
	if len(tokens) < 2 {
		return req
	}

	if tokens[0] == "Bearer" {
		// Treat the authorization header as session key bearer token
		userID, sessionCreationTime := t.onAuth(req.Context(), tokens[1])
		session.UserID = userID
		session.Creation = sessionCreationTime
		session.ShieldClientRole = auth.GQLShieldClientRegular
	} else if tokens[0] == "Debug" {
		// Treat the authorization header as debug session key bearer token
		if t.onDebugAuth(req.Context(), tokens[1]) {
			session.IsDebug = true
			session.ShieldClientRole = auth.GQLShieldClientDebug
		}
	}

	// Put context in the request context
	return req
}
