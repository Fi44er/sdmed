package utils

func FindDifferences[T any](oldObjs, newObjs []T, getName func(T) string) (deleted []string, added []string) {
	oldNames := make(map[string]bool)

	for _, obj := range oldObjs {
		oldNames[getName(obj)] = true
	}

	for _, obj := range newObjs {
		name := getName(obj)
		if oldNames[name] {
			delete(oldNames, name)
		} else {
			added = append(added, name)
		}
	}

	for name := range oldNames {
		deleted = append(deleted, name)
	}

	return deleted, added
}
