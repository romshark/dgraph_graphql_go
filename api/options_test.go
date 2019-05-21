package api_test

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api"
	"github.com/romshark/dgraph_graphql_go/api/transport"
	thttp "github.com/romshark/dgraph_graphql_go/api/transport/http"
	"github.com/stretchr/testify/require"
)

func TestOptionsInvalid(t *testing.T) {
	assumeErr := func(t *testing.T, opts api.ServerOptions) {
		require.Error(t, opts.Prepare())
	}

	t.Run("noTransport", func(t *testing.T) {
		assumeErr(t, api.ServerOptions{
			Mode:      api.ModeProduction,
			Transport: []transport.Server{},
		})
	})

	t.Run("production/debugUserEnabled", func(t *testing.T) {
		debugUsrOptions := []api.DebugUserStatus{
			api.DebugUserRW,
			api.DebugUserReadOnly,
		}
		for _, debugUsrOption := range debugUsrOptions {
			t.Run(string(debugUsrOption), func(t *testing.T) {
				serverHTTP, err := thttp.NewServer(thttp.ServerOptions{
					Host: "localhost:80",
					TLS: &thttp.ServerTLS{
						CertificateFilePath: "certfile",
						PrivateKeyFilePath:  "privkeyfile",
					},
					Playground: true,
				})
				require.NoError(t, err)
				require.NotNil(t, serverHTTP)

				assumeErr(t, api.ServerOptions{
					Mode:      api.ModeProduction,
					Transport: []transport.Server{serverHTTP},
					DebugUser: api.DebugUserOptions{
						Status: debugUsrOption,
					},
				})
			})
		}
	})

	t.Run("production/nonTLSTransport", func(t *testing.T) {
		serverHTTP, err := thttp.NewServer(thttp.ServerOptions{
			Host:       "localhost:80",
			TLS:        nil,
			Playground: true,
		})
		require.NoError(t, err)
		require.NotNil(t, serverHTTP)

		assumeErr(t, api.ServerOptions{
			Mode:      api.ModeProduction,
			Transport: []transport.Server{serverHTTP},
		})
	})
}
