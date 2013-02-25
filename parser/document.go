package parser

import (
	"andyk/docs/util"
	"regexp"
	"strings"
)

type DocumentParser struct {
	Patterns DocumentStructure
}

func NewDocumentParser(documentStructure DocumentStructure) DocumentParser {
	return DocumentParser{
		Patterns: documentStructure,
	}
}

func (parser DocumentParser) Parse(lines []string, metaData MetaData) (item ParsedItem, err error) {

	// assign meta data
	item.MetaData = metaData

	// title
	title, lines := parser.getMatchingValue(lines, parser.Patterns.Title)
	item.AddElement("title", title)

	// description
	description, lines := parser.getMatchingValue(lines, parser.Patterns.Description)
	item.AddElement("description", description)

	// content
	item.AddElement("content", parser.getContent(lines))

	return item, nil
}

func (parser DocumentParser) getMatchingValue(lines []string, pattern regexp.Regexp) (string, []string) {

	// In order to be the "matching value" the line must
	// either be empty or match the supplied pattern.
	for lineNumber, line := range lines {

		lineMatchesTitlePattern, matches := util.IsMatch(line, parser.Patterns.Title)
		if lineMatchesTitlePattern {
			nextLine := getNextLinenumber(lineNumber, lines)
			return util.GetLastElement(matches), lines[nextLine:]
		}

		lineIsEmpty := parser.Patterns.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return "", lines
}

func (parser DocumentParser) getContent(lines []string) string {

	startLine := 0
	endLine := len(lines)

	return strings.TrimSpace(strings.Join(lines[startLine:endLine], "\n"))
}

func getNextLinenumber(lineNumber int, lines []string) int {
	nextLine := lineNumber + 1

	if nextLine <= len(lines) {
		return nextLine
	}

	return lineNumber
}
