// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package presentation

import (
	"github.com/elWyatt/allmark/model"
	"github.com/elWyatt/allmark/services/parser/document"
	"github.com/elWyatt/allmark/services/parser/pattern"
	"fmt"
	"strings"
	"time"
)

func Parse(item *model.Item, lastModifiedDate time.Time, lines []string) (parseError error) {

	// parse as document
	if _, err := document.Parse(item, lastModifiedDate, lines); err != nil {
		return fmt.Errorf("Unable to parse item %q. Error: %s", item, err)
	}

	// convert the document to a presentation
	item.Content = convertToPresentation(item.Content)

	return
}

func convertToPresentation(content string) string {

	// split the lines again
	presentationLines := make([]string, 0)
	lines := strings.Split(content, "\n")

	// separate the slides with horizontal rule
	for lineNumber, line := range lines {

		// check if the current line is a headline
		isHeadline, headlineText, _ := pattern.IsHeadline(line)

		// skip all lines which are not headlines
		if !isHeadline {
			presentationLines = append(presentationLines, line)
			continue
		}

		// prepend a horizontal rule if
		// - its not the first line
		// - the headline is not already preceeded with a horizontal rule
		if lineNumber > 0 && !horizontalRuleAlreadyPresentIn(lines[0:lineNumber-1]) {

			presentationLines = append(presentationLines, "")
			presentationLines = append(presentationLines, "---")
			presentationLines = append(presentationLines, "")

		}

		// Fix the headline levels:
		// If the current line is followed by content make the current headline a level-two headline.
		// If the current line is not followed by content make the current headline a level-one headline.
		if lineNumber < len(lines)-1 && followingLinesContainContent(lines[lineNumber+1:]) {

			// slide with content -> h2 headline
			secondLevelHeadline := fmt.Sprintf("## %s", headlineText)
			presentationLines = append(presentationLines, secondLevelHeadline)

		} else {

			// slide without content -> h1 headline
			firstLevelHeadline := fmt.Sprintf("# %s", headlineText)
			presentationLines = append(presentationLines, firstLevelHeadline)

		}
	}

	// save the presentation code
	presentationContent := strings.TrimSpace(strings.Join(presentationLines, "\n"))

	return presentationContent
}

// Determine whether the supplied lines contain a horizontal rule
// before a line contains actual content.
func horizontalRuleAlreadyPresentIn(lines []string) bool {
	if len(lines) == 0 {
		return false
	}

	for lineNumber := len(lines) - 1; lineNumber >= 0; lineNumber-- {
		line := lines[lineNumber]

		// skip emppty lines
		if pattern.IsEmpty(line) {
			continue
		}

		return pattern.IsHorizontalRule(line)
	}

	return false
}

// Determine if the supplied lines contain content before
// the next slide-end (horizontal rule or headine).
func followingLinesContainContent(lines []string) bool {
	if len(lines) == 0 {
		return false
	}

	for _, line := range lines {

		// an empty line is not content.
		if pattern.IsEmpty(line) {
			continue
		}

		// if there is another headline, there is no more content.
		if isHeadline, _, _ := pattern.IsHeadline(line); isHeadline {
			return false
		}

		// if there is a horizontal rule, there is no more content.
		if pattern.IsHorizontalRule(line) {
			return false
		}

		// if it is not white-space, a headline or a
		// horizontal rule it must be content.
		return true
	}

	return false
}
