package gqlshield

import (
	"fmt"

	"github.com/pkg/errors"
)

func copyRoles(roles ...ClientRole) (map[int]ClientRole, error) {
	if len(roles) < 1 {
		return nil, errors.New("missing roles")
	}

	byName := make(map[string]ClientRole, len(roles))
	byID := make(map[int]ClientRole, len(roles))
	cp := make(map[int]ClientRole, len(roles))
	for _, role := range roles {
		// Ensure the client role ID is unique
		if defined, alreadyDefined := byID[role.ID]; alreadyDefined {
			return nil, fmt.Errorf(
				"role ID %d redefined for %s",
				defined.ID,
				defined.Name,
			)
		}

		// Ensure the client role name is valid
		if len(role.Name) < 1 {
			return nil, fmt.Errorf(
				"invalid (empty) client role name (%d)",
				role.ID,
			)
		}

		// Ensure the client role name is unique
		if defined, alreadyDefined := byName[role.Name]; alreadyDefined {
			return nil, fmt.Errorf(
				"%d:'%s' redefined for %d",
				defined.ID,
				role.Name,
				role.ID,
			)
		}

		role := ClientRole{
			ID:   role.ID,
			Name: role.Name,
		}

		// Define name->id
		byName[role.Name] = role
		byID[role.ID] = role
		cp[role.ID] = role
	}
	return cp, nil
}
