package dbmod

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

// ReactionSubject represents a union (Post | Reaction) type wrapper
type ReactionSubject struct {
	V interface{}
}

// UID returns the unique entity node ID
func (u ReactionSubject) UID() *string {
	switch v := u.V.(type) {
	case *Post:
		uid := v.UID
		return &uid
	case *Reaction:
		uid := v.UID
		return &uid
	}
	return nil
}

// MarshalJSON implements the Marshaler interface
func (u ReactionSubject) MarshalJSON() ([]byte, error) {
	switch v := u.V.(type) {
	case *Post:
		return json.Marshal(v)
	case *Reaction:
		return json.Marshal(v)
	}
	panic(errors.Errorf(
		"invalid union ReactionSubject value of type: %s",
		reflect.TypeOf(u.V),
	))
}

// UnmarshalJSON implements the Unmarshaler interface
func (u *ReactionSubject) UnmarshalJSON(d []byte) error {
	var keyVal map[string]interface{}
	if err := json.Unmarshal(d, &keyVal); err != nil {
		return err
	}
	if _, exists := keyVal["Reaction.id"]; exists {
		var v Reaction
		if err := json.Unmarshal(d, &v); err != nil {
			return err
		}
		u.V = &v
	} else if _, exists := keyVal["Post.id"]; exists {
		var v Post
		if err := json.Unmarshal(d, &v); err != nil {
			return err
		}
		u.V = &v
	} else {
		return errors.Errorf(
			"unsupported JSON for union ReactionSubject: '%s'",
			string(d),
		)
	}
	return nil
}
