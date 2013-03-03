package parser

import (
	"andyk/docs/indexer"
	"andyk/docs/util"
	"errors"
	"fmt"
	"github.com/russross/blackfriday"
	"os"
	"regexp"
	"strings"
)

var (
	// Lines which contain nothing but white space characters
	// or no characters at all.	
	EmptyLinePattern = regexp.MustCompile("^\\s*$")

	// Lines which a start with a hash, followed by zero or more
	// white space characters, followed by text.
	TitlePattern = regexp.MustCompile("\\s*#\\s*(\\w.+)")

	// Lines which start with text
	DescriptionPattern = regexp.MustCompile("^\\w.+")

	// Lines which nothing but dashes
	HorizontalRulePattern = regexp.MustCompile("^-{2,}")

	// Lines with a "key: value" syntax
	MetaDataPattern = regexp.MustCompile("^(\\w+):\\s*(\\w.+)$")
)

type ParsedItemElement struct {
	Name  string
	Value string
}

func NewParsedItemElement(name string, value string) ParsedItemElement {
	return ParsedItemElement{
		Name:  name,
		Value: value,
	}
}

type ParsedItem struct {
	Item     indexer.Item
	Elements []ParsedItemElement
	MetaData MetaData
}

func NewParsedItem(item indexer.Item) ParsedItem {
	return ParsedItem{
		Item: item,
	}
}

func (parsedItem *ParsedItem) GetElementValue(name string) string {
	if parsedItem.Elements == nil || len(parsedItem.Elements) == 0 {
		return ""
	}

	for _, element := range parsedItem.Elements {

		if strings.ToLower(element.Name) == strings.ToLower(name) {
			return element.Value
		}

	}

	return ""
}

func (parsedItem *ParsedItem) AddElement(name string, value string) {

	element := NewParsedItemElement(name, value)

	if parsedItem.Elements == nil {
		parsedItem.Elements = make([]ParsedItemElement, 1, 1)
		parsedItem.Elements[0] = element
		return
	}

	parsedItem.Elements = append(parsedItem.Elements, element)

}

func Parse(item indexer.Item) (ParsedItem, error) {

	// open the file
	file, err := os.Open(item.Path)
	if err != nil {
		return NewParsedItem(item), err
	}

	defer file.Close()

	// get the lines
	lines := util.GetLines(file)

	switch item.Type {
	case indexer.DocumentItemType:
		{
			parsedDoc := NewParsedItem(item)
			return *parsedDoc.ParseDocument(lines), nil
		}
	}

	return NewParsedItem(item), errors.New(fmt.Sprintf("Items of type \"%v\" cannot be parsed.", item.Type))
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

func getContent(lines []string) string {

	// all remaining lines are the (raw markdown) content
	rawMarkdownContent := strings.TrimSpace(strings.Join(lines, "\n"))

	// html encode the markdown
	htmlEncodedContent := string(blackfriday.MarkdownBasic([]byte(rawMarkdownContent)))

	return htmlEncodedContent
}

func getNextLinenumber(lineNumber int, lines []string) int {
	nextLine := lineNumber + 1

	if nextLine <= len(lines) {
		return nextLine
	}

	return lineNumber
}
