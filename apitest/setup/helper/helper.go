package helper

import (
	"testing"
	"time"

	"github.com/romshark/dgraph_graphql_go/api"
)

type successAssumption bool

const (
	success          successAssumption = true
	potentialFailure successAssumption = false
)

type testSetupInterface interface {
	T() *testing.T
	Query(q string, res interface{}) *api.ResponseError
	QueryVar(
		q string,
		vars map[string]string,
		res interface{},
	) *api.ResponseError
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
