package goconcat

// removes index and returns new slice
func RemoveFromSlice[T any](slice []T, index int) (poppedSlice []T) {
	for i, s := range slice {
		if i == index {
			continue
		}
		poppedSlice = append(poppedSlice, s)
	}
	return
}

// return index from slice
func PopFromSlice[T any](slice []T, index int) (singleElement T) {
	for i, s := range slice {
		if index == i {
			singleElement = s
		}
	}
	return
}

// return index from slice
func ReturnAllButIndices[T any](slices []T, indices []int) (newSlice []T) {
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

type StringConvert interface {
	~string
}

func AnyToString[T StringConvert](s T) string {
	return string(s)
}

type TSlice interface {
	~string
}

func StringToType[T TSlice](slice []string) []T {
	var newDirectories []T

	for _, s := range slice {
		newDirectories = append(newDirectories, T(s))
	}

	return newDirectories
}
