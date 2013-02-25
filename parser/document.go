package parser

import (
	"andyk/docs/util"
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
	title, lines := parser.getTitle(lines)
	item.AddElement("title", title)

	// description
	description, lines := parser.getDescription(lines)
	item.AddElement("description", description)

	// content
	item.AddElement("content", parser.getContent(lines))

	return item, nil
}

func (parser DocumentParser) getTitle(lines []string) (string, []string) {

	// In order to be the "title" the line must either
	// be empty or match the title pattern.
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

func (parser DocumentParser) getDescription(lines []string) (string, []string) {

	// In order to be a "description" the line must either
	// be empty or match the description pattern.
	for lineNumber, line := range lines {

		lineMatchesDescriptionPattern, matches := util.IsMatch(line, parser.Patterns.Description)
		if lineMatchesDescriptionPattern {
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
