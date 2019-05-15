package api

import "context"

// onRootAuth handles a root authentication request
func (srv *server) onRootAuth(
	ctx context.Context,
	sessionKey string,
) bool {
	return string(srv.rootSessionKey) == sessionKey
}
