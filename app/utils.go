package app

// map[namespace] -> map[name] -> exists
func selectNewNames(ol map[string]map[string]bool, nw map[string]map[string]bool) (result map[string]map[string]bool) {
	for ns, names := range nw {
		if _, exists := ol[ns]; !exists {
			result[ns] = names
			continue
		}

		for name, exists := range names {
			if !exists {
				continue
			}

			if _, inOld := ol[ns][name]; !inOld {
				if _, ok := result[ns]; !ok {
					result[ns] = make(map[string]bool, 100)
				}
				result[ns][name] = true
			}
		}
	}

	return
}
