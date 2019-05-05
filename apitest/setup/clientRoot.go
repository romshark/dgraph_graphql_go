package setup

import "github.com/stretchr/testify/require"

// Root creates a new authenticated API root client
func (ts *TestSetup) Root() *Client {
	clt := ts.Guest()

	// Sign in as root
	require.NoError(ts.t, clt.apiClient.AuthRoot(
		ts.rootUsername,
		ts.rootPassword,
	))

	return clt
}
