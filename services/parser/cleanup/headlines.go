// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cleanup

import (
	"regexp"
)

var (
	// markdown headline which start with the headline text righ after the hash
	markdownHeadlineStartWhitespace = regexp.MustCompile(`^(#+)([^\s#])`)

	// match headlines which have a hash at the end
	markdownHeadlineClosingHeadline = regexp.MustCompile(`^(#+)\s*(.+?)\s*(#+)$`)
)

func cleanupHeadlines(lines []string) []string {
	for index, line := range lines {

		fixedLine := line

		// headline start
		fixedLine = markdownHeadlineStartWhitespace.ReplaceAllString(fixedLine, "$1 $2")

		// remove closing headline hashes
		fixedLine = markdownHeadlineClosingHeadline.ReplaceAllString(fixedLine, "$1 $2")

		// same the fixed line
		if fixedLine != line {
			lines[index] = fixedLine
		}
	}

	return lines
}
