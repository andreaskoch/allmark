package util

// GetLastElement retrn the last element of a string array.
func GetLastElement(array []string) string {
	if array == nil {
		return ""
	}

	return array[len(array)-1]
}
