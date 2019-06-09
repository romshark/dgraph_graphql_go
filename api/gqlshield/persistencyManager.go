package gqlshield

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

// PersistencyManager represents a persistency manager
type PersistencyManager interface {
	// Load loads the GraphQL shield configuration
	Load() (*State, error)

	// Save persists the GraphQL shield configuration
	Save(*State) error
}

// NewPepersistencyManagerFileJSON creates a new JSON file based
// persistency manager
func NewPepersistencyManagerFileJSON(
	path string,
	syncWrite bool,
) (PersistencyManager, error) {
	return &persistencyManagerFileJSON{
		path:      path,
		syncWrite: syncWrite,
	}, nil
}

type persistencyManagerFileJSON struct {
	path      string
	syncWrite bool
}

func (man *persistencyManagerFileJSON) onFile(f func(*os.File) error) error {
	// Open file
	flags := os.O_CREATE | os.O_RDWR
	if man.syncWrite {
		flags = os.O_CREATE | os.O_RDWR | os.O_SYNC
	}
	file, err := os.OpenFile(man.path, flags, 0660)
	if err != nil {
		return errors.Wrap(err, "opening file")
	}
	defer file.Close()

	return f(file)
}

func (man *persistencyManagerFileJSON) Load() (state *State, err error) {
	err = man.onFile(func(file *os.File) error {
		jsonDecoder := json.NewDecoder(file)
		state = &State{}
		return jsonDecoder.Decode(state)
	})
	return
}

func (man *persistencyManagerFileJSON) Save(state *State) error {
	return man.onFile(func(file *os.File) error {
		jsonEncoder := json.NewEncoder(file)
		state = &State{}
		return jsonEncoder.Encode(state)
	})
}
