package auth

// IsUser indicates that the client is required to be an authenticated user
type IsUser struct{}

func (rule IsUser) check(session *RequestSession) string {
	if session.UserID == "" {
		return "the user is required to be an authenticated user"
	}
	return ""
}
