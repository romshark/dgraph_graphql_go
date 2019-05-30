package gqlshield

import "errors"

func prepareQuery(query []byte) ([]byte, error) {
	if len(query) < 1 {
		return nil, errors.New("invalid (empty) query")
	}

	start := int(-1)
	shift := int(0)
	tail := len(query)
	inString := false

	// shift over leading spaces
	i := 0
	for ; i < len(query); i++ {
		char := query[i]
		if char == ' ' || char == '\t' || char == '\n' {
		} else {
			break
		}
	}
	tail -= i

	for ; i < len(query); i++ {
		char := query[i]
		if char == ' ' || char == '\t' || char == '\n' {
			if !inString && start < 0 {
				// record spaces start
				start = shift
			}
		} else if start > -1 {
			// shift over spaces
			query[start] = ' '
			delta := shift - start
			if delta > 1 {
				tail -= delta - 1
			}
			shift = start + 1
			start = -1
		}
		if char == '"' {
			inString = !inString
		}
		query[shift] = char
		shift++
	}
	if start > -1 {
		tail -= shift - start
	}
	if inString {
		return nil, errors.New("unclosed string context")
	}
	return query[:tail], nil
}
