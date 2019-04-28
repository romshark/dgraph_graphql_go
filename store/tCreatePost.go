package store

import (
	"context"
)

// CreatePost creates a new post
func (str *store) CreatePost(
	ctx context.Context,
	author ID,
	title string,
	contents string,
) (err error) {
	// Validate input
	if err := ValidatePostTitle(title); err != nil {
		return err
	}
	if err := ValidatePostContents(contents); err != nil {
		return err
	}

	//TODO: implement
	/*
		// Prepare
		newPostID := NewID()

		// Begin transaction
		txn, close := str.newTxn(&err)
		defer close()
	*/

	return
}
