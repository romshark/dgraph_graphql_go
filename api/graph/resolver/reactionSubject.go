package resolver

// ReactionSubject implements the identically named union type
type ReactionSubject struct {
	subject interface{}
}

// ToReaction casts the union to a Reaction resolver
func (un *ReactionSubject) ToReaction() (*Reaction, bool) {
	res, ok := un.subject.(*Reaction)
	return res, ok
}

// ToPost casts the union to a Post resolver
func (un *ReactionSubject) ToPost() (*Post, bool) {
	res, ok := un.subject.(*Post)
	return res, ok
}
