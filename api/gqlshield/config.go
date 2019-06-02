package gqlshield

// Config defines the GraphQL shield configuration
type Config struct {
	// PersistencyManager is used for configuration state persistency.
	// Persistency is disabled if PersistencyManager is nil.
	PersistencyManager PersistencyManager
}
