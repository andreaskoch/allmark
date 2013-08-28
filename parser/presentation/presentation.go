// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package presentation

import (
	"fmt"
	"github.com/andreaskoch/allmark/parser/document"
	"github.com/andreaskoch/allmark/parser/pattern"
	"github.com/andreaskoch/allmark/repository"
	"regexp"
	"strings"
)

var (
	// markdown headline which start with the headline text righ after the hash
	markdownHeadlineStartWhitespace = regexp.MustCompile(`^(#+)([^\s#])`)
)

// Parse an item with a title, description and content
func Parse(item *repository.Item, lines []string, fallbackTitle string) (sucess bool, err error) {

	// parse the document
	if success, err := document.Parse(item, lines, fallbackTitle); !success {
		return sucess, err
	}

	// split the lines again
	presentationLines := make([]string, 0)
	lines = strings.Split(item.RawContent, "\n")

	highestHeadlineLevel := 6
	lowestHeadlineLevel := 1

	// separate the slides with horizontal rule
	for lineNumber, line := range lines {

		// skip non-headlines
		if !isHeadline(line) {
			presentationLines = append(presentationLines, line)
			continue
		}

		// determine the headline level
		headlineLevel, err := getHeadlineLevel(line)
		if err != nil {
			panic(err)
		}

		// capture the highest headline level
		if headlineLevel < highestHeadlineLevel {
			highestHeadlineLevel = headlineLevel
		}

		// capture the lowest headline level
		if headlineLevel > lowestHeadlineLevel {
			lowestHeadlineLevel = headlineLevel
		}

		// get the lines before this line
		if lineNumber > 0 && horizontalRuleAlreadyPresentIn(lines[0:lineNumber-1]) {

			// headline is already preceeded by horizontal rule
			presentationLines = append(presentationLines, line)

		} else {

			// prepend a horizontal rule
			presentationLines = append(presentationLines, "")
			presentationLines = append(presentationLines, "---")
			presentationLines = append(presentationLines, "")

			presentationLines = append(presentationLines, line)
		}

	}

	// normalize the headline levels
	// for lineNumber, line := range lines {

	// }

	// save the presentation code
	item.RawContent = strings.TrimSpace(strings.Join(presentationLines, "\n"))

	return true, nil
}

func getHeadlineLevel(line string) (int, error) {

	if !isHeadline(line) {
		return 0, fmt.Errorf("The line %q is not a headline.", line)
	}

	level := 0
	for _, character := range line {

		if string(character) == `#` {
			level++
		} else {
			break
		}
	}

	return level, nil

}

func isHeadline(line string) bool {
	return strings.HasPrefix(line, "#")
}

func horizontalRuleAlreadyPresentIn(lines []string) bool {
	if len(lines) == 0 {
		return false
	}

	for lineNumber := len(lines) - 1; lineNumber >= 0; lineNumber-- {
		line := lines[lineNumber]

		if pattern.EmptyLinePattern.MatchString(line) {
			continue
		}

		return pattern.HorizontalRulePattern.MatchString(line)
	}

	return false
}
