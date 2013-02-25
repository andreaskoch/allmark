package parser

import (
	"andyk/docs/util"
	"errors"
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
	titleLocation, lines := parser.locateTitle(lines)
	if !titleLocation.Found {
		return item, errors.New("Title not found.")
	}

	item.AddElement("title", getTitle(titleLocation))

	// description
	descriptionLocation, lines := parser.locateDescription(lines)
	if !descriptionLocation.Found {
		return item, errors.New("Description not found.")
	}

	item.AddElement("description", getDescription(descriptionLocation))

	// content
	contentLocation := parser.locateContent(lines)
	if !contentLocation.Found {
		return item, errors.New("No content available.")
	}

	item.AddElement("description", getContent(contentLocation))

	return item, nil
}

func getTitle(titleLocation Match) string {
	return strings.TrimSpace(util.GetLastElement(titleLocation.Matches))
}

func getDescription(descriptionLocation Match) string {
	return strings.TrimSpace(util.GetLastElement(descriptionLocation.Matches))
}

func getContent(contentLocation Match) string {
	return strings.TrimSpace(strings.Join(contentLocation.Matches, "\n"))
}

func (parser DocumentParser) locateTitle(lines []string) (Match, []string) {

	// In order to be the "title" the line must either
	// be empty or match the title pattern.
	for lineNumber, line := range lines {

		lineMatchesTitlePattern, matches := util.IsMatch(line, parser.Patterns.Title)
		if lineMatchesTitlePattern {
			nextLine := getNextLinenumber(lineNumber, lines)
			return Found(matches), lines[nextLine:]
		}

		lineIsEmpty := parser.Patterns.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return NotFound(), lines
}

func getNextLinenumber(lineNumber int, lines []string) int {
	nextLine := lineNumber + 1

	if nextLine <= len(lines) {
		return nextLine
	}

	return lineNumber
}

func (parser DocumentParser) locateDescription(lines []string) (Match, []string) {

	// In order to be a "description" the line must either
	// be empty or match the description pattern.
	for lineNumber, line := range lines {

		lineMatchesDescriptionPattern, matches := util.IsMatch(line, parser.Patterns.Description)
		if lineMatchesDescriptionPattern {
			nextLine := getNextLinenumber(lineNumber, lines)
			return Found(matches), lines[nextLine:]
		}

		lineIsEmpty := parser.Patterns.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return NotFound(), lines
}

func (parser DocumentParser) locateContent(lines []string) Match {

	startLine := 0
	endLine := len(lines)

	// All lines between the start- and endLine are content
	return Found(lines[startLine:endLine])
}
