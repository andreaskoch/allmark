// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"
	"github.com/andreaskoch/allmark/markdown"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/types"
	"github.com/andreaskoch/allmark/util"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Lines which contain nothing but white space characters
	// or no characters at all.
	EmptyLinePattern = regexp.MustCompile(`^\s*$`)

	// Lines which a start with a hash, followed by zero or more
	// white space characters, followed by text.
	TitlePattern = regexp.MustCompile(`^#\s*([\pL\pN\p{Latin}]+.+)`)

	// Lines which start with text
	DescriptionPattern = regexp.MustCompile(`^[\pL\pN\p{Latin}]+.+`)

	// Lines which nothing but dashes
	HorizontalRulePattern = regexp.MustCompile(`^-{2,}`)

	// Lines with a "key: value" syntax
	MetaDataPattern = regexp.MustCompile(`^(\w+):\s*([\pL\pN\p{Latin}]+.+)$`)
)

func Parse(item *repository.Item) (*repository.Item, error) {
	if item.IsVirtual() {
		return parseVirtual(item)
	}

	return parsePhysical(item)
}

func parseVirtual(item *repository.Item) (*repository.Item, error) {

	if item == nil {
		return nil, fmt.Errorf("Cannot create meta data from nil.")
	}

	// get the item title
	title := filepath.Base(item.Directory())

	// create the meta data
	metaData, err := newMetaData(item)
	if err != nil {
		return nil, err
	}

	item.Title = title
	item.MetaData = metaData

	return item, nil
}

func parsePhysical(item *repository.Item) (*repository.Item, error) {

	// open the file
	file, err := os.Open(item.Path())
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	defer file.Close()

	// get the raw lines
	lines := util.GetLines(file)

	// parse the meta data
	fallbackItemTypeFunc := func() string {
		return getItemTypeFromFilename(item.Path())
	}

	item.MetaData, lines = parseMetaData(item, lines, fallbackItemTypeFunc)

	// parse the content
	switch itemType := item.MetaData.ItemType; itemType {
	case types.DocumentItemType, types.CollectionItemType, types.RepositoryItemType, types.PresentationItemType:
		{
			if success, err := parseDocumentLikeItem(item, lines); success {
				return item, nil
			} else {
				return nil, err
			}
		}
	case types.MessageItemType:
		{
			if success, err := parseMessage(item, lines); success {
				return item, nil
			} else {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("Item %q (type: %s) cannot be parsed.", item.Path(), itemType)
	}

	panic("Unreachable")
}

// Parse an item with a title, description and content
func parseDocumentLikeItem(item *repository.Item, lines []string) (sucess bool, err error) {

	// title
	item.Title, lines = getTitle(lines)

	// description
	item.Description, lines = getDescription(lines)

	// raw markdown content
	item.RawContent = lines

	return true, nil
}

func parseMessage(item *repository.Item, lines []string) (sucess bool, err error) {

	// raw markdown content
	item.RawContent = lines

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

func getTitle(lines []string) (title string, remainingLines []string) {
	title, remainingLines = getMatchingValue(lines, TitlePattern)

	// cleanup the title
	title = strings.TrimSpace(title)
	title = strings.TrimSuffix(title, "#")

	return title, remainingLines
}

func getDescription(lines []string) (description string, remainingLines []string) {
	description, remainingLines = getMatchingValue(lines, DescriptionPattern)

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

func getItemTypeFromFilename(filenameOrPath string) string {

	if !markdown.IsMarkdownFile(filenameOrPath) {
		return types.UnknownItemType // abort if file does not have a markdown extension
	}

	extension := filepath.Ext(filenameOrPath)
	filenameWithExtension := filepath.Base(filenameOrPath)
	filename := filenameWithExtension[0:(strings.LastIndex(filenameWithExtension, extension))]

	switch strings.ToLower(filename) {
	case types.DocumentItemType:
		return types.DocumentItemType

	case types.PresentationItemType:
		return types.PresentationItemType

	case types.CollectionItemType:
		return types.CollectionItemType

	case types.MessageItemType:
		return types.MessageItemType

	case types.RepositoryItemType:
		return types.RepositoryItemType

	default:
		return types.DocumentItemType
	}

	return types.UnknownItemType
}
