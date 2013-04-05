// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

// GetLastElement retrn the last element of a string array.
func GetLastElement(array []string) string {
	if array == nil {
		return ""
	}

	return array[len(array)-1]
}
