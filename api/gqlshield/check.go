package gqlshield

import "fmt"

func (shld *shield) Check(
	clientRoleID int,
	queryString []byte,
	arguments map[string]string,
) ([]byte, error) {
	if len(queryString) < 1 {
		return queryString, Error{
			Code:    ErrWrongInput,
			Message: "invalid (empty) query",
		}
	}

	normalized, err := prepareQuery(queryString)
	if err != nil {
		return queryString, err
	}

	if shld.conf.WhitelistOption != WhitelistEnabled {
		// Don't check the query if query whitelisting is disabled
		return normalized, nil
	}

	shld.lock.RLock()
	defer shld.lock.RUnlock()

	// Find role
	if _, roleDefined := shld.clientRoles[clientRoleID]; !roleDefined {
		return normalized, fmt.Errorf("role %d is undefined", clientRoleID)
	}

	// Lookup query
	qrObj, found := shld.index.Search(normalized)
	if !found {
		return normalized, Error{
			Code:    ErrUnauthorized,
			Message: "query not whitelisted",
		}
	}
	qr := qrObj.(*query)

	// Ensure the client is allowed to execute this query
	if _, roleAllowed := qr.whitelistedFor[clientRoleID]; !roleAllowed {
		return normalized, Error{
			Code: ErrUnauthorized,
			Message: fmt.Sprintf(
				"role %d is not allowed to execute this query",
				clientRoleID,
			),
		}
	}

	// Check arguments
	if len(arguments) != len(qr.parameters) {
		return normalized, Error{
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
			return normalized, Error{
				Code:    ErrUnauthorized,
				Message: fmt.Sprintf("missing argument '%s'", name),
			}
		}
		if uint32(len(actual)) > expectedParam.MaxValueLength {
			return normalized, Error{
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

	return normalized, nil
}
