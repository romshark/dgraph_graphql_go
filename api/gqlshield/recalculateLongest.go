package gqlshield

func (shld *shield) recalculateLongest() error {
	// Recalculate longest
	shld.longest = 0
	for itr := shld.index.Iterator(); itr.HasNext(); {
		node, err := itr.Next()
		if err != nil {
			return err
		}
		queryLength := len(node.Value().(*query).query)
		if queryLength > shld.longest {
			shld.longest = queryLength
		}
	}
	return nil
}
