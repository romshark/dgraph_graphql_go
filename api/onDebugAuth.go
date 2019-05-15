package api

import "context"

// onDebugAuth handles a debug client authentication request
func (srv *server) onDebugAuth(
	ctx context.Context,
	sessionKey string,
) bool {
	return string(srv.debugSessionKey) == sessionKey
}
