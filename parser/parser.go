// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	UnknownItemType      = "unknown"
	DocumentItemType     = "document"
	PresentationItemType = "presentation"
	CollectionItemType   = "collection"
	MessageItemType      = "message"
	RepositoryItemType   = "repository"
)

var (
	// Lines which contain nothing but white space characters
	// or no characters at all.
	EmptyLinePattern = regexp.MustCompile(`^\s*$`)

	// Lines which a start with a hash, followed by zero or more
	// white space characters, followed by text.
	TitlePattern = regexp.MustCompile(`^#\s*(\w.+)`)

	// Lines which start with text
	DescriptionPattern = regexp.MustCompile(`^\w.+`)

	// Lines which nothing but dashes
	HorizontalRulePattern = regexp.MustCompile(`^-{2,}`)

	// Lines with a "key: value" syntax
	MetaDataPattern = regexp.MustCompile(`^(\w+):\s*(\w.+)$`)
)

type ParsedItem struct {
	*repository.Item

	Title       string
	Description string
	RawContent  []string
	MetaData    MetaData

	ConvertedContent string
}

func Parse(lines []string, item *repository.Item) (*ParsedItem, error) {

	// parse meta data
	result := &ParsedItem{
		Item: item,
	}

	result.MetaData, lines = parseMetaData(lines, func() string {
		return getItemTypeFromFilename(item.Path())
	})

	itemType := result.MetaData.ItemType
	switch itemType {
	case DocumentItemType, CollectionItemType, RepositoryItemType:
		{
			if success, err := parseDocumentLikeItem(result, lines); success {
				return result, nil
			} else {
				return nil, err
			}
		}
	case MessageItemType:
		{
			if success, err := parseMessage(result, lines); success {
				return result, nil
			} else {
				return nil, err
			}
		}
	}

	return nil, fmt.Errorf("Item %q (type: %s) cannot be parsed.", item.Path(), itemType)
}

// Parse an item with a title, description and content
func parseDocumentLikeItem(parserParsedItem *ParsedItem, lines []string) (sucess bool, err error) {

	// title
	parserParsedItem.Title, lines = getMatchingValue(lines, TitlePattern)

	// description
	parserParsedItem.Description, lines = getMatchingValue(lines, DescriptionPattern)

	// raw markdown content
	parserParsedItem.RawContent = lines

	return true, nil
}

func parseMessage(parserParsedItem *ParsedItem, lines []string) (sucess bool, err error) {

	// raw markdown content
	parserParsedItem.RawContent = lines

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

func getItemTypeFromFilename(filenameOrPath string) string {
	extension := strings.ToLower(filepath.Ext(filenameOrPath))

	if extension != ".md" && extension != ".mdown" && extension != ".markdown" {
		return UnknownItemType // abort if file does not have a markdown extension
	}

	filenameWithExtension := filepath.Base(filenameOrPath)
	filename := filenameWithExtension[0:(strings.LastIndex(filenameWithExtension, extension))]

	switch strings.ToLower(filename) {
	case DocumentItemType:
		return DocumentItemType

	case PresentationItemType:
		return PresentationItemType

	case CollectionItemType:
		return CollectionItemType

	case MessageItemType:
		return MessageItemType

	case RepositoryItemType:
		return RepositoryItemType

	default:
		return DocumentItemType
	}

	return UnknownItemType
}
