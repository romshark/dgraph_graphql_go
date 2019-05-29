package config

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/romshark/dgraph_graphql_go/api/passhash"
	"github.com/romshark/dgraph_graphql_go/api/sesskeygen"
	thttp "github.com/romshark/dgraph_graphql_go/api/transport/http"
)

// File represents a TOML encoded configuration file
type File struct {
	Mode                Mode                `toml:"mode"`
	PasswordHasher      PasswordHasher      `toml:"password-hasher"`
	SessionKeyGenerator SessionKeyGenerator `toml:"session-key-generator"`
	DB                  struct {
		Host string `toml:"host"`
	} `toml:"db"`
	Log struct {
		Debug string `toml:"debug"`
		Error string `toml:"error"`
	} `toml:"log"`
	Debug struct {
		Mode     string `toml:"mode"`
		Username string `toml:"username"`
		Password string `toml:"password"`
	} `toml:"debug"`
	TransportHTTP struct {
		Host              string   `toml:"host"`
		KeepAliveDuration Duration `toml:"keep-alive-duration"`
		Playground        bool     `toml:"playground"`
		TLS               struct {
			Enabled          bool             `toml:"enabled"`
			MinVersion       TLSVersion       `toml:"min-version"`
			CertificateFile  string           `toml:"certificate-file"`
			KeyFile          string           `toml:"key-file"`
			CurvePreferences []TLSCurveID     `toml:"curve-preferences"`
			CipherSuites     []TLSCipherSuite `toml:"cipher-suites"`
		} `toml:"tls"`
	} `toml:"transport-http"`
}

func (f *File) mode(conf *ServerConfig) error {
	if err := f.Mode.Validate(); err != nil {
		return err
	}
	conf.Mode = f.Mode
	return nil
}

func (f *File) dbHost(conf *ServerConfig) error {
	conf.DBHost = f.DB.Host
	return nil
}

func (f *File) passwordHasher(conf *ServerConfig) error {
	switch f.PasswordHasher {
	case "bcrypt":
		conf.PasswordHasher = passhash.Bcrypt{}
		return nil
	}
	return fmt.Errorf("unsupported password hasher: '%s'", f.PasswordHasher)
}

func (f *File) sessionKeyGenerator(conf *ServerConfig) error {
	switch f.SessionKeyGenerator {
	case "default":
		conf.SessionKeyGenerator = sesskeygen.NewDefault()
		return nil
	}
	return fmt.Errorf(
		"unsupported session key generator: '%s'",
		f.SessionKeyGenerator,
	)
}

func (f *File) debugLog(conf *ServerConfig) error {
	var writer io.Writer
	if strings.HasPrefix(f.Log.Debug, "stdout") {
		writer = os.Stdout
	} else if strings.HasPrefix(f.Log.Debug, "file:") &&
		len(f.Log.Debug) > 5 {
		// Debug log to file
		var err error
		writer, err = os.OpenFile(
			f.Log.Debug,
			os.O_WRONLY|os.O_APPEND|os.O_CREATE,
			0660,
		)
		if err != nil {
			return errors.Wrap(err, "debug log file")
		}
	} else {
		return fmt.Errorf("invalid: '%s'", f.Log.Debug)
	}
	conf.DebugLog = log.New(
		writer,
		"DBG: ",
		log.Ldate|log.Ltime,
	)
	return nil
}

func (f *File) errorLog(conf *ServerConfig) error {
	var writer io.Writer
	if strings.HasPrefix(f.Log.Error, "stderr") {
		writer = os.Stdout
	} else if strings.HasPrefix(f.Log.Error, "file:") &&
		len(f.Log.Error) > 5 {
		// Error log to file
		var err error
		writer, err = os.OpenFile(
			f.Log.Error,
			os.O_WRONLY|os.O_APPEND|os.O_CREATE,
			0660,
		)
		if err != nil {
			return errors.Wrap(err, "error log file")
		}
	} else {
		return fmt.Errorf("invalid: '%s'", f.Log.Error)
	}
	conf.ErrorLog = log.New(
		writer,
		"ERR: ",
		log.Ldate|log.Ltime,
	)
	return nil
}

func (f *File) debug(conf *ServerConfig) error {
	conf.DebugUser.Mode = DebugUserMode(f.Debug.Mode)
	if err := conf.DebugUser.Mode.Validate(); err != nil {
		return errors.Wrap(err, "debug user mode")
	}
	conf.DebugUser.Username = f.Debug.Username
	conf.DebugUser.Password = f.Debug.Password
	return nil
}

func (f *File) transportHTTP(conf *ServerConfig) error {
	opt := thttp.ServerConfig{}

	// Host
	if len(f.TransportHTTP.Host) < 1 {
		return nil
	}

	// Keep-alive duration
	opt.KeepAliveDuration = time.Duration(f.TransportHTTP.KeepAliveDuration)

	// TLS
	if f.TransportHTTP.TLS.Enabled {
		opt.TLS = &thttp.ServerTLS{
			Config:              &tls.Config{},
			CertificateFilePath: f.TransportHTTP.TLS.CertificateFile,
			PrivateKeyFilePath:  f.TransportHTTP.TLS.KeyFile,
		}

		// Min version
		opt.TLS.Config.MinVersion = uint16(f.TransportHTTP.TLS.MinVersion)

		// Cipher suites
		cipherSuites := make([]uint16, len(f.TransportHTTP.TLS.CipherSuites))
		for i, cipherSuite := range f.TransportHTTP.TLS.CipherSuites {
			cipherSuites[i] = uint16(cipherSuite)
		}
		opt.TLS.Config.CipherSuites = cipherSuites

		// Curve preferences
		curveIDs := make([]tls.CurveID, len(f.TransportHTTP.TLS.CipherSuites))
		for i, curveID := range f.TransportHTTP.TLS.CipherSuites {
			curveIDs[i] = tls.CurveID(curveID)
		}
		opt.TLS.Config.CurvePreferences = curveIDs
	}

	// Playground
	opt.Playground = f.TransportHTTP.Playground

	newServer, err := thttp.NewServer(opt)
	if err != nil {
		return errors.Wrap(err, "HTTP server init")
	}

	conf.Transport = append(conf.Transport, newServer)

	return nil
}

// FromFile reads the configuration from a file
func FromFile(path string) (*ServerConfig, error) {
	var file File
	conf := &ServerConfig{}

	// Read TOML config file
	if _, err := toml.DecodeFile(path, &file); err != nil {
		return nil, errors.Wrap(err, "TOML decode")
	}

	for setterName, setter := range map[string]func(*ServerConfig) error{
		"mode":                  file.mode,
		"db.host":               file.dbHost,
		"password-hasher":       file.passwordHasher,
		"session-key-generator": file.sessionKeyGenerator,
		"log.debug":             file.debugLog,
		"log.error":             file.errorLog,
		"debug":                 file.debug,
		"transport-http":        file.transportHTTP,
	} {
		if err := setter(conf); err != nil {
			return nil, errors.Wrap(err, setterName)
		}
	}

	if err := conf.Prepare(); err != nil {
		return nil, err
	}

	return conf, nil
}
