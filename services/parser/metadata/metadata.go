// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metadata

import (
	"github.com/andreaskoch/allmark2/services/parser/pattern"
)

func GetLines(lines []string) []string {

	if len(lines) == 0 {
		return []string{}
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
				return lines[lineNumber+1:]
			}

			// no meta data detected
			return []string{}
		}
	}

	// no meta data detected
	return []string{}
}
