package util

func MergeStringMap(base map[string]string, toMerge map[string]string) map[string]string {
	for k, v := range toMerge {
		if _, ok := base[k]; !ok {
			base[k] = v
		}
	}

	return base
}
