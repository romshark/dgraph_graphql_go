package helper

import (
	"testing"
	"time"
)

type successAssumption bool

const (
	success          successAssumption = true
	potentialFailure successAssumption = false
)

type testSetupInterface interface {
	T() *testing.T
	Query(q string, res interface{}) []string
	QueryVar(q string, vars map[string]string, res interface{}) []string
}

// New creates a new test helper
func New(testSetup testSetupInterface) Helper {
	helper := Helper{
		ts:                     testSetup,
		creationTimeTollerance: time.Second * 3,
	}
	helper.OK = AssumeSuccess{
		h: &helper,
		t: testSetup.T(),
	}
	return helper
}

// Helper represents a test helper
type Helper struct {
	ts                     testSetupInterface
	creationTimeTollerance time.Duration
	OK                     AssumeSuccess
}

// AssumeSuccess wraps failable helper functions
type AssumeSuccess struct {
	h *Helper
	t *testing.T
}
