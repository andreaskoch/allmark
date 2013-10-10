// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

// GetLastElement retrn the last element of a string array.
func GetLastElement(slice []string) string {
	if slice == nil {
		return ""
	}

	return slice[len(slice)-1]
}

func SliceContainsElement(list []string, elem string) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}
	return false
}

func TrimSlice(slice []string) []string {
	trimmed := make([]string, 0)

	for _, element := range slice {
		if element == "" {
			continue
		}

		trimmed = append(trimmed, element)
	}

	return trimmed
}
