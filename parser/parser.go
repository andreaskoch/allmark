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
	TitlePattern = regexp.MustCompile(`^#\s*(\pL.+)`)

	// Lines which start with text
	DescriptionPattern = regexp.MustCompile(`^\pL.+`)

	// Lines which nothing but dashes
	HorizontalRulePattern = regexp.MustCompile(`^-{2,}`)

	// Lines with a "key: value" syntax
	MetaDataPattern = regexp.MustCompile(`^(\w+):\s*(\pN.+)$`)
)

type ParsedItem struct {
	*repository.Item

	Title       string
	Description string
	RawContent  []string
	MetaData    MetaData

	ConvertedContent string
}

func Parse(item *repository.Item) (*ParsedItem, error) {
	if item.IsVirtual() {
		return parseVirtual(item)
	}

	return parsePhysical(item)
}

func parseVirtual(item *repository.Item) (*ParsedItem, error) {

	if item == nil {
		return nil, fmt.Errorf("Cannot create meta data from nil.")
	}

	title := filepath.Base(item.Directory())

	metaData, err := newMetaData(item)
	if err != nil {
		return nil, err
	}

	result := &ParsedItem{
		Item: item,

		Title:    title,
		MetaData: metaData,
	}

	return result, nil
}

func parsePhysical(item *repository.Item) (*ParsedItem, error) {

	// open the file
	file, err := os.Open(item.Path())
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	defer file.Close()

	// get the raw lines
	lines := util.GetLines(file)

	// create a result
	result := &ParsedItem{
		Item: item,
	}

	// parse the meta data
	fallbackItemTypeFunc := func() string {
		return getItemTypeFromFilename(item.Path())
	}

	result.MetaData, lines = parseMetaData(item, lines, fallbackItemTypeFunc)

	// parse the content
	switch itemType := result.MetaData.ItemType; itemType {
	case types.DocumentItemType, types.CollectionItemType, types.RepositoryItemType:
		{
			if success, err := parseDocumentLikeItem(result, lines); success {
				return result, nil
			} else {
				return nil, err
			}
		}
	case types.MessageItemType:
		{
			if success, err := parseMessage(result, lines); success {
				return result, nil
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
