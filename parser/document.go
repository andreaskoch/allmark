package parser

import (
	"andyk/docs/util"
	"strings"
)

var documentStructure = NewDocumentStructure()

type Document struct {
	Title       string
	Description string
	Content     string
	MetaData    MetaData
}

// CreateDocument returns a new Document from the given Item.
func CreateDocument(lines []string) Document {
	doc := Document{}

	// assign meta data
	metaData, metaDataLocation := GetMetaData(lines, documentStructure)
	if !metaDataLocation.Found {
		return doc
	}
	doc.MetaData = metaData

	// assign title
	titleLocation := locateTitle(lines)
	if !titleLocation.Found {
		return doc
	}
	doc.Title = getTitle(titleLocation)

	// description
	descriptionLocation := locateDescription(lines, titleLocation)
	if !descriptionLocation.Found {
		return doc
	}
	doc.Description = getDescription(descriptionLocation)

	// content
	contentLocation := locateContent(lines, descriptionLocation, metaDataLocation)
	if !contentLocation.Found {
		return doc
	}
	doc.Content = getContent(contentLocation)

	return doc
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

func locateTitle(lines []string) Match {

	// In order to be the "title" the line must either
	// be empty or match the title pattern.

	for lineNumber, line := range lines {

		lineMatchesTitlePattern, matches := util.IsMatch(line, documentStructure.Title)
		if lineMatchesTitlePattern {
			return Found(lineNumber, lineNumber, matches)
		}

		lineIsEmpty := documentStructure.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return NotFound()
}

func locateDescription(lines []string, titleLocation Match) Match {

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

		lineMatchesDescriptionPattern, matches := util.IsMatch(line, documentStructure.Description)
		if lineMatchesDescriptionPattern {
			absoluteLineNumber := startLine + relativeLineNumber
			return Found(absoluteLineNumber, absoluteLineNumber, matches)
		}

		lineIsEmpty := documentStructure.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return NotFound()
}

func locateContent(lines []string, descriptionLocation Match, metaDataLocation Match) Match {

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
	if metaDataLocation.Found {
		endLine = metaDataLocation.Lines.Start - 1
	} else {
		endLine = len(lines)
	}

	// All lines between the start- and endLine are content
	return Found(startLine, endLine, lines[startLine:endLine])
}
