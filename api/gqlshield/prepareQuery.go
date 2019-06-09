package gqlshield

func prepareQuery(query []byte) ([]byte, error) {
	if len(query) < 1 {
		return nil, Error{
			Code:    ErrWrongInput,
			Message: "invalid (empty) query",
		}
	}

	start := int(-1)
	shift := int(0)
	tail := len(query)
	inString := false

	// shift over leading spaces
	i := 0
LEADING_LOOP:
	for ; i < len(query); i++ {
		char := query[i]
		if char == '\\' && i+1 < len(query) {
			switch query[i+1] {
			case 't':
				// escaped tab
				fallthrough
			case 'n':
				// escaped line-break
				fallthrough
			case 'r':
				// escaped carriage return
				i++
			default:
				break LEADING_LOOP
			}
		} else if char == ' ' || char == '\t' || char == '\n' {
		} else {
			break
		}
	}
	tail -= i

	for ; i < len(query); i++ {
		char := query[i]
		if char == '\\' && i+1 < len(query) {
			switch query[i+1] {
			case 't':
				// escaped tab
				fallthrough
			case 'n':
				// escaped line-break
				fallthrough
			case 'r':
				// escaped carriage return
				if start < 0 {
					start = shift
				}
				i++
				tail--
			default:
			}
		} else if !inString && (char == ' ' || char == '\t' || char == '\n') {
			if start < 0 {
				// record shift start
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
		return nil, Error{
			Code:    ErrWrongInput,
			Message: "unclosed string context",
		}
	}
	return query[:tail], nil
}
