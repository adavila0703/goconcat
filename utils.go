package goconcat

// removes index and returns new slice
func removeFromSlice[T any](slice []T, index int) (poppedSlice []T) {
	for i, s := range slice {
		if i == index {
			continue
		}
		poppedSlice = append(poppedSlice, s)
	}
	return
}

// return index from slice
func popFromSlice[T any](slice []T, index int) (singleElement T) {
	for i, s := range slice {
		if index == i {
			singleElement = s
		}
	}
	return
}

// return index from slice
func returnAllButIndices[T any](slices []T, indices []int) (newSlice []T) {
	ignoredIndex := make(map[int]int)

	// map slice of ints
	for _, num := range indices {
		ignoredIndex[num] = num
	}

	for index, slice := range slices {
		if _, ok := ignoredIndex[index]; ok {
			continue
		}

		newSlice = append(newSlice, slice)
	}
	return
}

type stringConvert interface {
	~string
}

func anyToString[T stringConvert](s T) string {
	return string(s)
}

type tSlice interface {
	~string
}

func stringToType[T tSlice](slice []string) []T {
	var newDirectories []T

	for _, s := range slice {
		newDirectories = append(newDirectories, T(s))
	}

	return newDirectories
}
