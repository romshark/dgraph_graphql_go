package emotion

import "github.com/pkg/errors"

// Emotion represents an emotion type
type Emotion string

const (
	// Happy represents happiness
	Happy Emotion = "happy"

	// Angry represents anger
	Angry Emotion = "angry"

	// Excited represents excitement
	Excited Emotion = "excited"

	// Fearful represents anxiety
	Fearful Emotion = "fearful"

	// Thoughtful represents thoughtfulness
	Thoughtful Emotion = "thoughtful"
)

// Validate returns an error if the value is invalid
func Validate(v Emotion) error {
	switch v {
	case Happy:
		fallthrough
	case Angry:
		fallthrough
	case Excited:
		fallthrough
	case Fearful:
		fallthrough
	case Thoughtful:
		return nil
	}
	return errors.Errorf("invalid value: '%s'", v)
}
