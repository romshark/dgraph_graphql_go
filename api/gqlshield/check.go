package gqlshield

import "fmt"

func (shld *shield) Check(
	clientRoleID int,
	Query []byte,
	arguments map[string]string,
) error {
	if len(Query) < 1 {
		return Error{
			Code:    ErrWrongInput,
			Message: "invalid (empty) query",
		}
	}
	normalized, err := prepareQuery(Query)
	if err != nil {
		return err
	}

	shld.lock.RLock()
	defer shld.lock.RUnlock()

	// Find role
	if _, roleDefined := shld.clientRoles[clientRoleID]; !roleDefined {
		return fmt.Errorf("role %d is undefined", clientRoleID)
	}

	// Lookup query
	qrObj, found := shld.index.Search(normalized)
	if !found {
		return Error{
			Code:    ErrUnauthorized,
			Message: "query not whitelisted",
		}
	}
	qr := qrObj.(*query)

	// Ensure the client is allowed to execute this query
	if _, roleAllowed := qr.whitelistedFor[clientRoleID]; !roleAllowed {
		return Error{
			Code: ErrUnauthorized,
			Message: fmt.Sprintf(
				"role %d is not allowed to execute this query",
				clientRoleID,
			),
		}
	}

	// Check arguments
	if len(arguments) != len(qr.parameters) {
		return Error{
			Code: ErrUnauthorized,
			Message: fmt.Sprintf(
				"unexpected number of arguments: (%d/%d)",
				len(arguments),
				len(qr.parameters),
			),
		}
	}
	for name, expectedParam := range qr.parameters {
		actual, hasArg := arguments[name]
		if !hasArg {
			return Error{
				Code:    ErrUnauthorized,
				Message: fmt.Sprintf("missing argument '%s'", name),
			}
		}
		if uint32(len(actual)) > expectedParam.MaxValueLength {
			return Error{
				Code: ErrUnauthorized,
				Message: fmt.Sprintf(
					"argument '%s' exceeds max length (%d/%d)",
					name,
					len(actual),
					expectedParam.MaxValueLength,
				),
			}
		}
	}

	return nil
}
