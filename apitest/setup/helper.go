package setup

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/romshark/dgraph_graphql_go/store/errors"
)

// Helper represents a test helper
type Helper struct {
	ts                     *TestSetup
	c                      *Client
	creationTimeTollerance time.Duration
	OK                     AssumeSuccess
	ERR                    AssumeFailure
}

// AssumeSuccess wraps failable helper functions
type AssumeSuccess struct {
	h *Helper
	t *testing.T
}

// AssumeFailure wraps failable helper functions
type AssumeFailure struct {
	h *Helper
	t *testing.T
}

func (notOk AssumeFailure) checkErrCode(errCode errors.Code) {
	require.NotEqual(notOk.t, "", errors.FilterCode(errCode))
}
