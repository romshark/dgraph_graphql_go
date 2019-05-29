package config

import "fmt"

// DebugUserMode represents a debug user mode
type DebugUserMode string

const (
	// DebugUserUnset represents the default unset option value
	DebugUserUnset DebugUserMode = ""

	// DebugUserDisabled disables the debug user
	DebugUserDisabled DebugUserMode = "disabled"

	// DebugUserReadOnly enables the debug user in a read-only mode
	DebugUserReadOnly DebugUserMode = "read-only"

	// DebugUserRW enables the debug user in a read-write mode
	DebugUserRW DebugUserMode = "read-write"
)

// Validate returns an error if the value is invalid
func (md DebugUserMode) Validate() error {
	switch md {
	case DebugUserUnset:
		fallthrough
	case DebugUserDisabled:
		fallthrough
	case DebugUserRW:
		fallthrough
	case DebugUserReadOnly:
		return nil
	}
	return fmt.Errorf("invalid debug user mode: '%s'", md)
}
