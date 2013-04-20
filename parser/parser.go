// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"regexp"
)

var (
	// Lines which contain nothing but white space characters
	// or no characters at all.	
	EmptyLinePattern = regexp.MustCompile(`^\s*$`)

	// Lines which a start with a hash, followed by zero or more
	// white space characters, followed by text.
	TitlePattern = regexp.MustCompile(`\s*#\s*(\w.+)`)

	// Lines which start with text
	DescriptionPattern = regexp.MustCompile(`^\w.+`)

	// Lines which nothing but dashes
	HorizontalRulePattern = regexp.MustCompile(`^-{2,}`)

	// Lines with a "key: value" syntax
	MetaDataPattern = regexp.MustCompile(`^(\w+):\s*(\w.+)$`)
)

type Result struct {
	Title       string
	Description string
	RawContent  []string
	Type        string
	MetaData    repository.MetaData

	ConvertedContent string
}

func Parse(lines []string, itemType string) (*Result, error) {

	// parse meta data
	result := &Result{}
	result.MetaData, lines = parseMetaData(lines)

	switch itemType {
	case repository.DocumentItemType, repository.CollectionItemType, repository.RepositoryItemType:
		{
			if success, err := parseDocumentLikeItem(result, lines); success {
				return result, nil
			} else {
				return nil, err
			}
		}
	case repository.MessageItemType:
		{
			if success, err := parseMessage(result, lines); success {
				return result, nil
			} else {
				return nil, err
			}
		}
	}

	return nil, fmt.Errorf("Items of type \"%v\" cannot be parsed.", itemType)
}

// Parse an item with a title, description and content
func parseDocumentLikeItem(parserResult *Result, lines []string) (sucess bool, err error) {

	// title
	parserResult.Title, lines = getMatchingValue(lines, TitlePattern)

	// description
	parserResult.Description, lines = getMatchingValue(lines, DescriptionPattern)

	// raw markdown content
	parserResult.RawContent = lines

	return true, nil
}

func parseMessage(parserResult *Result, lines []string) (sucess bool, err error) {

	// raw markdown content
	parserResult.RawContent = lines

	return true, nil
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

		lineIsEmpty := EmptyLinePattern.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return "", lines
}

func getNextLinenumber(lineNumber int, lines []string) int {
	nextLine := lineNumber + 1

	if nextLine <= len(lines) {
		return nextLine
	}

	return lineNumber
}
