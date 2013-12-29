// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package files

func trimSlice(slice []string) []string {
	trimmed := make([]string, 0)

	for _, element := range slice {
		if element == "" {
			continue
		}

		trimmed = append(trimmed, element)
	}

	return trimmed
}
