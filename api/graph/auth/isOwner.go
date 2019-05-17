package auth

import (
	"github.com/romshark/dgraph_graphql_go/store"
)

// IsOwner indicates that the client is required to be the owner of a resource
type IsOwner struct {
	Owner store.ID
}

func (rule IsOwner) check(session *RequestSession) string {
	if session.UserID != rule.Owner {
		return "the user is required to be the resource owner"
	}
	return ""
}
