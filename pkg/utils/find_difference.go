package utils

func FindDifferences[T any](oldObjs, newObjs []T, getInfo func(T) (string, string)) (deleted []string, added []string) {
	oldNames := make(map[string]bool)

	for _, obj := range oldObjs {
		name, _ := getInfo(obj)
		oldNames[name] = true
	}

	for _, obj := range newObjs {
		name, id := getInfo(obj)
		if oldNames[name] {
			delete(oldNames, id)
		} else {
			added = append(added, name)
		}
	}

	for id := range oldNames {
		deleted = append(deleted, id)
	}

	return deleted, added
}
