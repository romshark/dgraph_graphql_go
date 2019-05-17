package setup

import (
	"testing"
	"time"
)

type successAssumption bool

const (
	success          successAssumption = true
	potentialFailure successAssumption = false
)

// Helper represents a test helper
type Helper struct {
	ts                     *TestSetup
	c                      *Client
	creationTimeTollerance time.Duration
	OK                     AssumeSuccess
}

// AssumeSuccess wraps failable helper functions
type AssumeSuccess struct {
	h *Helper
	t *testing.T
}
