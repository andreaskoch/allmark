// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package document

import (
	"github.com/andreaskoch/allmark/parser/pattern"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"regexp"
	"strings"
)

var (
	// markdown headline which start with the headline text righ after the hash
	markdownHeadlineStartWhitespace = regexp.MustCompile(`^(#+)([\S])`)

	// match headlines which have a hash at the end
	markdownHeadlineClosingHeadline = regexp.MustCompile(`^(#+)\s+(.+?)\s+(#+)$`)
)

// Parse an item with a title, description and content
func Parse(item *repository.Item, lines []string, fallbackTitle string) (sucess bool, err error) {

	// title
	item.Title, lines = getTitle(lines)
	if item.Title == "" {
		item.Title = fallbackTitle
	}

	// description
	item.Description, lines = getDescription(lines)

	// raw markdown content
	item.RawContent = strings.TrimSpace(strings.Join(cleanMarkdown(lines), "\n"))

	return true, nil
}

func cleanMarkdown(lines []string) []string {
	for index, line := range lines {

		fixedLine := line

		// headline start
		fixedLine = markdownHeadlineStartWhitespace.ReplaceAllString(fixedLine, "$1 $2")

		// remove closing headline hashes
		fixedLine = markdownHeadlineClosingHeadline.ReplaceAllString(fixedLine, "$2 $3")

		// same the fixed line
		if fixedLine != line {
			lines[index] = fixedLine
		}
	}

	return lines
}

func getMatchingValue(lines []string, matchPattern *regexp.Regexp) (string, []string) {

	// In order to be the "matching value" the line must
	// either be empty or match the supplied pattern.
	for lineNumber, line := range lines {

		lineMatchesTitlePattern, matches := util.IsMatch(line, matchPattern)
		if lineMatchesTitlePattern {
			nextLine := getNextLinenumber(lineNumber, lines)
			return util.GetLastElement(matches), lines[nextLine:]
		}

		lineIsEmpty := pattern.EmptyLinePattern.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return "", lines
}

func getTitle(lines []string) (title string, remainingLines []string) {
	title, remainingLines = getMatchingValue(lines, pattern.TitlePattern)

	// cleanup the title
	title = strings.TrimSpace(title)
	title = strings.TrimSuffix(title, "#")

	return title, remainingLines
}

func getDescription(lines []string) (description string, remainingLines []string) {
	description, remainingLines = getMatchingValue(lines, pattern.DescriptionPattern)

	// cleanup the description
	description = strings.TrimSpace(description)

	return description, remainingLines
}

func getNextLinenumber(lineNumber int, lines []string) int {
	nextLine := lineNumber + 1

	if nextLine <= len(lines) {
		return nextLine
	}

	return lineNumber
}
