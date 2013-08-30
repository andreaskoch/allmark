// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package presentation

import (
	"github.com/andreaskoch/allmark/parser/document"
	"github.com/andreaskoch/allmark/parser/pattern"
	"github.com/andreaskoch/allmark/repository"
	"regexp"
	"strings"
)

var (
	// markdown headline pattern
	anyLevelMarkdownHeadline = regexp.MustCompile(`^(#+?)([^#].+?[^#])(#*)$`)
)

// Parse an item with a title, description and content
func Parse(item *repository.Item, lines []string, fallbackTitle string) (sucess bool, err error) {

	// use the document parser. a presentation has the same structure as a document
	if success, err := document.Parse(item, lines, fallbackTitle); !success {
		return sucess, err
	}

	// split the lines again
	presentationLines := make([]string, 0)
	lines = strings.Split(item.RawContent, "\n")

	// separate the slides with horizontal rule
	for lineNumber, line := range lines {

		// skip non-headlines
		if !isHeadline(line) {
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
			secondLevelHeadline := anyLevelMarkdownHeadline.ReplaceAllString(line, "## $2")
			presentationLines = append(presentationLines, secondLevelHeadline)

		} else {

			// slide without content -> h1 headline
			firstLevelHeadline := anyLevelMarkdownHeadline.ReplaceAllString(line, "# $2")
			presentationLines = append(presentationLines, firstLevelHeadline)

		}
	}

	// save the presentation code
	item.RawContent = strings.TrimSpace(strings.Join(presentationLines, "\n"))

	return true, nil
}

// Determine whether the supplied text
// is a markdown headline.
func isHeadline(text string) bool {
	return strings.HasPrefix(text, "#")
}

// Determine whether the supplied lines contain a horizontal rule
// before a line contains actual content.
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

	panic("Unreachable")
}

// Determine if the supplied lines contain content before
// the next slide-end (horizontal rule or headine).
func followingLinesContainContent(lines []string) bool {
	if len(lines) == 0 {
		return false
	}

	for _, line := range lines {

		// an empty line is not content.
		if pattern.EmptyLinePattern.MatchString(line) {
			continue
		}

		// if there is another headline, there is no more content.
		if isHeadline(line) {
			return false
		}

		// if there is a horizontal rule, there is no more content.
		if pattern.HorizontalRulePattern.MatchString(line) {
			return false
		}

		// if it is not white-space, a headline or a
		// horizontal rule it must be content.
		return true
	}

	panic("Unreachable")
}
