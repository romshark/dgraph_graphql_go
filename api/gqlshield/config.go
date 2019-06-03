package gqlshield

// WhitelistOption represents the query whitelist option
type WhitelistOption byte

const (
	_ WhitelistOption = iota

	// WhitelistDisabled disables query whitelisting
	WhitelistDisabled

	// WhitelistEnabled enables query whitelisting
	WhitelistEnabled
)

// Config defines the GraphQL shield configuration
type Config struct {
	// WhitelistOption enables query whitelisting when true
	WhitelistOption WhitelistOption

	// PersistencyManager is used for configuration state persistency.
	// Persistency is disabled if PersistencyManager is nil
	PersistencyManager PersistencyManager
}

// SetDefaults sets the default configuration options
func (conf *Config) SetDefaults() {
	if conf.WhitelistOption == 0 {
		// Enable query whitelisting by default
		conf.WhitelistOption = WhitelistEnabled
	}
}
