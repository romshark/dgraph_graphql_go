package api

// onRootAuth handles a root authentication request
func (srv *server) onRootAuth(username, password string) ([]byte, bool) {
	// Check root credentials
	if username != srv.opts.RootUser.Username ||
		password != srv.opts.RootUser.Password {
		return nil, false
	}

	// Return session key
	return srv.rootSessionKey, true
}
