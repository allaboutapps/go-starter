package util

// ContainsString checks whether the given string slice contains the string provided.
func ContainsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}

	return false
}

// ContainsAllString checks whether the given string slice contains all strings provided.
func ContainsAllString(slice []string, sub ...string) bool {
	contains := make(map[string]bool)
	for _, v := range sub {
		contains[v] = false
	}

	for _, v := range slice {
		if _, ok := contains[v]; ok {
			contains[v] = true
		}
	}

	for _, v := range contains {
		if !v {
			return false
		}
	}

	return true
}

// UniqueString takes the string slice provided and returns a new slice with all duplicates removed.
func UniqueString(slice []string) []string {
	seen := make(map[string]struct{})
	res := make([]string, 0)

	for _, s := range slice {
		if _, ok := seen[s]; !ok {
			res = append(res, s)
			seen[s] = struct{}{}
		}
	}

	return res
}
