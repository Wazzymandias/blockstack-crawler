package names

// Names is a string slice
type Names []string

// NamespaceNames stores a namespace with its list of names
type NamespaceNames struct {
	Namespace string
	Names
}

// SelectNew takes two sets of names and returns all elements that exist in the second but not in first.
func SelectNew(ol map[string]map[string]bool, nw map[string]map[string]bool) map[string]map[string]bool {
	result := make(map[string]map[string]bool)

	for ns, names := range nw {
		if _, exists := ol[ns]; !exists {
			result[ns] = names
			continue
		}

		result[ns] = make(map[string]bool)

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

	return result
}

// MapToSlice converts the boolean map of names associated with a namespace to a string slice
// This allows for a more JSON friendly struct for marshalling.
func MapToSlice(in map[string]map[string]bool) (out map[string][]string) {
	out = make(map[string][]string)

	for k, v := range in {
		out[k] = make([]string, 1)

		for e := range v {
			out[k] = append(out[k], e)
		}
	}

	return
}
