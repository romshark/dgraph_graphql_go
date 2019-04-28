package store

import "context"

// CreateReaction creates a new post
func (str *store) CreateReaction(
	ctx context.Context,
	post ID,
	author ID,
	message string,
) (err error) {
	// Validate inputs
	if err := ValidateReactionMessage(message); err != nil {
		return err
	}

	//TODO: implement
	/*
		// Prepare
		newReactionID := NewID()

		// Begin transaction
		txn, close := str.newTxn(&err)
		defer close()
	*/

	return
}
