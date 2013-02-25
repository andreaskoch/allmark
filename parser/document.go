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
	// locate title
	titleLocation := parser.locateTitle(lines)
	if !titleLocation.Found {
		return item, errors.New("Title not found.")
	}

	// save title
	item.AddElement("title", getTitle(titleLocation))

	// exclude title from lines
	lines = lines[titleLocation.Lines.End:]

	// description
	// locate the description
	descriptionLocation := parser.locateDescription(lines)
	if !descriptionLocation.Found {
		return item, errors.New("Description not found.")
	}

	// save the description
	item.AddElement("description", getDescription(descriptionLocation))

	// exclude description from lines
	lines = lines[descriptionLocation.Lines.End:]

	// content
	// locate the content
	contentLocation := parser.locateContent(lines)
	if !contentLocation.Found {
		return item, errors.New("No content available.")
	}

	// save the content
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

func (parser DocumentParser) locateTitle(lines []string) Match {

	// In order to be the "title" the line must either
	// be empty or match the title pattern.
	for lineNumber, line := range lines {

		lineMatchesTitlePattern, matches := util.IsMatch(line, parser.Patterns.Title)
		if lineMatchesTitlePattern {
			return Found(lineNumber, lineNumber, matches)
		}

		lineIsEmpty := parser.Patterns.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return NotFound()
}

func (parser DocumentParser) locateDescription(lines []string) Match {

	// In order to be a "description" the line must either
	// be empty or match the description pattern.
	for lineNumber, line := range lines {

		lineMatchesDescriptionPattern, matches := util.IsMatch(line, parser.Patterns.Description)
		if lineMatchesDescriptionPattern {
			return Found(lineNumber, lineNumber, matches)
		}

		lineIsEmpty := parser.Patterns.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return NotFound()
}

func (parser DocumentParser) locateContent(lines []string) Match {

	startLine := 0
	endLine := len(lines)

	// All lines between the start- and endLine are content
	return Found(startLine, endLine, lines[startLine:endLine])
}
