package parser

import (
	"andyk/docs/util"
	"errors"
	"strings"
)

type DocumentParser struct {
	Patterns         DocumentStructure
	MetaData         MetaData
	MetaDataLocation Match
}

func NewDocumentParser(documentStructure DocumentStructure, metaData MetaData, metaDataLocation Match) DocumentParser {
	return DocumentParser{
		Patterns:         documentStructure,
		MetaData:         metaData,
		MetaDataLocation: metaDataLocation,
	}
}

func (parser DocumentParser) Parse(lines []string) (item ParsedItem, err error) {

	// assign meta data
	item.MetaData = parser.MetaData

	// assign title
	titleLocation := parser.locateTitle(lines)
	if !titleLocation.Found {
		return item, errors.New("Title not found.")
	}
	item.AddElement("title", getTitle(titleLocation))

	// description
	descriptionLocation := parser.locateDescription(lines, titleLocation)
	if !descriptionLocation.Found {
		return item, errors.New("Description not found.")
	}
	item.AddElement("description", getDescription(descriptionLocation))

	// content
	contentLocation := parser.locateContent(lines, descriptionLocation)
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

func (parser DocumentParser) locateDescription(lines []string, titleLocation Match) Match {

	// The description must be preceeded by a title
	if !titleLocation.Found {
		return NotFound()
	}

	// If the document has no more lines than the line
	// in which the title has been located, there
	// will be no room for a description
	startLine := titleLocation.Lines.End + 1
	if len(lines) <= startLine {
		return NotFound()
	}

	// In order to be a "description" the line must either
	// be empty or match the description pattern.
	for relativeLineNumber, line := range lines[startLine:] {

		lineMatchesDescriptionPattern, matches := util.IsMatch(line, parser.Patterns.Description)
		if lineMatchesDescriptionPattern {
			absoluteLineNumber := startLine + relativeLineNumber
			return Found(absoluteLineNumber, absoluteLineNumber, matches)
		}

		lineIsEmpty := parser.Patterns.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return NotFound()
}

func (parser DocumentParser) locateContent(lines []string, descriptionLocation Match) Match {

	// Content must be preceeded by a description
	if !descriptionLocation.Found {
		return NotFound()
	}

	// If the document has no more lines than the line
	// in which the description has been located, there
	// will be no room for content
	startLine := descriptionLocation.Lines.End + 1
	if len(lines) <= startLine {
		return NotFound()
	}

	// If the document contains meta data
	// the content will be between the description
	// and the meta data. If not the content
	// will go up to the end of the document.
	endLine := 0
	if parser.MetaDataLocation.Found {
		endLine = parser.MetaDataLocation.Lines.Start - 1
	} else {
		endLine = len(lines)
	}

	// All lines between the start- and endLine are content
	return Found(startLine, endLine, lines[startLine:endLine])
}
