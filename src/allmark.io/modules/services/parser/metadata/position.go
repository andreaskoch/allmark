// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"allmark.io/modules/services/parser/pattern"
	"fmt"
)

// Get the position of the meta data section from the supplied lines.
func GetMetaDataPosition(lines []string) (int, error) {

	if len(lines) == 0 {
		return 0, fmt.Errorf("There cannot be any meta data if the supplied lines are empty.")
	}

	hasMetaDataDefinition := false
	for lineNumber := len(lines) - 1; lineNumber >= 0; lineNumber-- {
		line := lines[lineNumber]

		// skip empty lines
		if pattern.IsEmpty(line) {
			continue
		}

		// check if a line contains a meta data definition
		if !hasMetaDataDefinition && pattern.IsMetaDataDefinition(line) {
			hasMetaDataDefinition = true
			continue
		}

		// abort if a horizontal rule has been found
		if pattern.IsHorizontalRule(line) {
			if hasMetaDataDefinition {
				return lineNumber, nil
			}

			// no meta data detected
			return 0, fmt.Errorf("No meta data found.")
		}
	}

	// no meta data detected
	return 0, fmt.Errorf("No meta data found.")
}
