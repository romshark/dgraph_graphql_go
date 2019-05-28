package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"

	"github.com/romshark/dgraph_graphql_go/api"
	"github.com/romshark/dgraph_graphql_go/api/options"
	"github.com/romshark/dgraph_graphql_go/api/passhash"
	"github.com/romshark/dgraph_graphql_go/api/sesskeygen"
	"github.com/romshark/dgraph_graphql_go/api/transport"
	thttp "github.com/romshark/dgraph_graphql_go/api/transport/http"
)

var host = flag.String("host", "localhost:16000", "API server host address")
var dbHost = flag.String("dbhost", "localhost:9080", "database host address")
var argCertFilePath = flag.String(
	"tlscert",
	"./demo.crt",
	"path to the TLS certificate file",
)
var argPrivateKeyFile = flag.String(
	"tlskey",
	"./demo.key",
	"path to the TLS private-key file",
)

func main() {
	flag.Parse()

	// Enable TLS if a certificate file is provided
	var tlsConf *thttp.ServerTLS
	if *argCertFilePath != "" {
		tlsConf = &thttp.ServerTLS{
			Config: &tls.Config{
				MinVersion: tls.VersionTLS12,
				CurvePreferences: []tls.CurveID{
					tls.X25519,
					tls.CurveP256,
				},
				PreferServerCipherSuites: true,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_AES_128_GCM_SHA256,
				},
			},
			CertificateFilePath: *argCertFilePath,
			PrivateKeyFilePath:  *argPrivateKeyFile,
		}
	}

	// Use HTTP as transport
	transportHTTP, err := thttp.NewServer(thttp.ServerOptions{
		Host:       *host,
		TLS:        tlsConf,
		Playground: true,
	})
	if err != nil {
		log.Fatalf("API server HTTP(S) transport init: %s", err)
	}

	api, err := api.NewServer(options.ServerOptions{
		Mode:                options.ModeBeta,
		Host:                *host,
		DBHost:              *dbHost,                 // database host address
		SessionKeyGenerator: sesskeygen.NewDefault(), // session key generator
		PasswordHasher:      passhash.Bcrypt{},       // password hasher
		Transport: []transport.Server{
			// HTTP(S) transport
			transportHTTP,
		},
	})
	if err != nil {
		log.Fatalf("API server init: %s", err)
	}

	if err := api.Launch(); err != nil {
		log.Fatalf("API server launch: %s", err)
	}

	// Setup termination signal listener
	onTerminate(func() {
		if err := api.Shutdown(context.Background()); err != nil {
			log.Fatalf("API server shutdown: %s", err)
		}
	})

	api.AwaitShutdown()
}
