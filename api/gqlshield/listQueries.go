package gqlshield

func (shld *shield) ListQueries() (map[string]Query, error) {
	shld.lock.RLock()
	defer shld.lock.RUnlock()

	allQueries := make(map[string]Query, shld.index.Size())
	for itr := shld.index.Iterator(); itr.HasNext(); {
		node, err := itr.Next()
		if err != nil {
			return nil, err
		}
		qr := node.Value().(*query)
		allQueries[qr.name] = qr
	}
	return allQueries, nil
}
