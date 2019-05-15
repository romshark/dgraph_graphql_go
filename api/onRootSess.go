package api

import "context"

// onRootSess handles a root authentication request
func (srv *server) onRootSess(
	ctx context.Context,
	username, password string,
) []byte {
	// Check root credentials
	if username != srv.opts.RootUser.Username ||
		password != srv.opts.RootUser.Password {
		return nil
	}

	// Return session key
	return srv.rootSessionKey
}
