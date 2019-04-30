package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/stretchr/testify/require"
)

// TestCreatePostErr tests all possible post creation errors
func TestCreatePostErr(t *testing.T) {
	// Test invalid title
	t.Run("invalidTitle", func(t *testing.T) {
		invalidTitles := map[string]string{
			"empty":    "",
			"tooShort": "f",
			"tooLong": "11111111000000001111111100000000" +
				"11111111000000001111111100000000" +
				"2",
		}

		for tName, invalidTitle := range invalidTitles {
			t.Run(tName, func(t *testing.T) {
				ts := setup.New(t, tcx)
				defer ts.Teardown()

				author := ts.Help.OK.CreateUser("fooBarowich", "foo@bar.buz")
				res, err := ts.Help.CreatePost(
					*author.ID,
					invalidTitle,
					"test contents",
				)
				require.Nil(t, res)
				verifyError(t, "InvalidInput", err)
			})
		}
	})
}
