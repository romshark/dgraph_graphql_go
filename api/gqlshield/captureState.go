package gqlshield

func (shld *shield) captureState() *State {
	roles := make([]ClientRole, 0, len(shld.clientRoles))
	for _, role := range shld.clientRoles {
		roles = append(roles, role)
	}

	queries := make(map[string]QueryModel, len(shld.queriesByName))
	for _, query := range shld.queriesByName {
		var params map[string]Parameter
		if query.parameters != nil {
			params = make(map[string]Parameter, len(query.parameters))
			for name, param := range query.parameters {
				params[name] = param
			}
		}

		whitelistedFor := make([]int, 0, len(query.whitelistedFor))
		for role := range query.whitelistedFor {
			whitelistedFor = append(whitelistedFor, role)
		}

		queries[string(query.id)] = QueryModel{
			Query:          string(query.query),
			Creation:       query.creation,
			Name:           query.name,
			Parameters:     params,
			WhitelistedFor: whitelistedFor,
		}
	}

	return &State{
		Roles:              roles,
		WhitelistedQueries: queries,
	}
}
