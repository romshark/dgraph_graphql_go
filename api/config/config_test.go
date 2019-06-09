package config_test

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/config"
	"github.com/romshark/dgraph_graphql_go/api/transport"
	thttp "github.com/romshark/dgraph_graphql_go/api/transport/http"
	"github.com/stretchr/testify/require"
)

func TestConfigInvalid(t *testing.T) {
	assumeErr := func(t *testing.T, conf config.ServerConfig) {
		require.Error(t, conf.Prepare())
	}

	t.Run("noTransport", func(t *testing.T) {
		assumeErr(t, config.ServerConfig{
			Mode:      config.ModeProduction,
			Transport: []transport.Server{},
		})
	})

	t.Run("production/debugUserEnabled", func(t *testing.T) {
		debugUserModes := []config.DebugUserMode{
			config.DebugUserRW,
			config.DebugUserReadOnly,
		}
		for _, debugUserMode := range debugUserModes {
			t.Run(string(debugUserMode), func(t *testing.T) {
				serverHTTP, err := thttp.NewServer(thttp.ServerConfig{
					Host: "localhost:80",
					TLS: &thttp.ServerTLS{
						CertificateFilePath: "certfile",
						PrivateKeyFilePath:  "privkeyfile",
					},
					Playground: true,
				})
				require.NoError(t, err)
				require.NotNil(t, serverHTTP)

				assumeErr(t, config.ServerConfig{
					Mode:      config.ModeProduction,
					Transport: []transport.Server{serverHTTP},
					DebugUser: config.DebugUserConfig{
						Mode: debugUserMode,
					},
				})
			})
		}
	})

	t.Run("production/nonTLSTransport", func(t *testing.T) {
		serverHTTP, err := thttp.NewServer(thttp.ServerConfig{
			Host:       "localhost:80",
			TLS:        nil,
			Playground: true,
		})
		require.NoError(t, err)
		require.NotNil(t, serverHTTP)

		assumeErr(t, config.ServerConfig{
			Mode:      config.ModeProduction,
			Transport: []transport.Server{serverHTTP},
		})
	})
}
