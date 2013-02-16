package model

import (
	"regexp"
)

type Document struct {
	Title       string
	Description string
	Content     string
	Hash        string

	pattern  DocumentPattern
	rawLines []string
}

func CreateDocument(repositoryItem *RepositoryItem) *Document {
	doc := Document{
		Hash:     repositoryItem.GetHash(),
		pattern:  NewDocumentPattern(),
		rawLines: repositoryItem.GetLines(),
	}

	// parse
	return doc.parse()
}

func (doc *Document) parse() *Document {
	return setTitle(doc)
}

type DocumentPattern struct {
	EmptyLine      regexp.Regexp
	Title          regexp.Regexp
	Description    regexp.Regexp
	HorizontalRule regexp.Regexp
	MetaData       regexp.Regexp
}

func NewDocumentPattern() DocumentPattern {
	emptyLineRegexp := regexp.MustCompile("^\\s*$")
	titleRegexp := regexp.MustCompile("\\s*#\\s*(.+)")
	descriptionRegexp := regexp.MustCompile("^\\w.+")
	horizontalRuleRegexp := regexp.MustCompile("^-{2,}")
	metaDataRegexp := regexp.MustCompile("^(\\w+):\\s*(\\w.+)$")

	return DocumentPattern{
		EmptyLine:      *emptyLineRegexp,
		Title:          *titleRegexp,
		Description:    *descriptionRegexp,
		HorizontalRule: *horizontalRuleRegexp,
		MetaData:       *metaDataRegexp,
	}
}

// Check if the current Document contains a title
func (doc *Document) locateTitle() (found bool, lineNumber int) {

	// In order to be the "title" the line must either
	// be empty or match the title pattern.

	for lineNumber, line := range doc.rawLines {

		lineMatchesTitlePattern := doc.pattern.Title.MatchString(line)
		if lineMatchesTitlePattern {
			return true, lineNumber
		}

		lineIsEmpty := doc.pattern.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return false, 0
}

// Check if the current Document contains a description
func (doc *Document) locateDescription() (found bool, lineNumber int) {

	// The description must be preceeded by a title
	titleExists, titleLineNumber := doc.locateTitle()
	if !titleExists {
		return false, 0
	}

	// If the document has no more lines than the line
	// in which the title has been located, there
	// will be no room for a description
	startLine := titleLineNumber + 1
	if len(doc.rawLines) <= startLine {
		return false, 0
	}

	// In order to be a "description" the line must either
	// be empty or match the description pattern.
	for lineNumber, line := range doc.rawLines[startLine:] {

		lineMatchesDescriptionPattern := doc.pattern.Description.MatchString(line)
		if lineMatchesDescriptionPattern {
			return true, lineNumber
		}

		lineIsEmpty := doc.pattern.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return false, 0
}

// Check if the current Document contains meta data
func (doc *Document) locateMetaData() (found bool, lineNumber int) {

	// Find the last horizontal rule in the document
	lastFoundHorizontalRulePosition := -1
	for lineNumber, line := range doc.rawLines {

		lineMatchesHorizontalRulePattern := doc.pattern.HorizontalRule.MatchString(line)
		if lineMatchesHorizontalRulePattern {
			lastFoundHorizontalRulePosition = lineNumber
		}

	}

	// If there is no horizontal rule there is no meta data
	if lastFoundHorizontalRulePosition == -1 {
		return false, 0
	}

	// If the document has no more lines than
	// the last found horizontal rule there is no
	// room for meta data
	metaDataStartLine := lastFoundHorizontalRulePosition + 1
	if len(doc.rawLines) <= metaDataStartLine {
		return false, 0
	}

	// Check if the last horizontal rule is followed
	// either by white space or be meta data
	for lineNumber, line := range doc.rawLines[metaDataStartLine:] {

		lineMatchesMetaDataPattern := doc.pattern.MetaData.MatchString(line)
		if lineMatchesMetaDataPattern {
			return true, metaDataStartLine
		}

		lineIsEmpty := doc.pattern.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			return false, 0
		}

	}

	return false, 0
}

// Check if the current Document contains content
func (doc *Document) locateContent() (found bool, startLine int, endLine int) {

	// Content must be preceeded by a description
	descriptionExists, descriptionLineNumber := doc.locateDescription()
	if !descriptionExists {
		return false, 0, 0
	}

	// If the document has no more lines than the line
	// in which the description has been located, there
	// will be no room for content
	startLine = descriptionLineNumber + 1
	if len(doc.rawLines) <= startLine {
		return false, 0, 0
	}

	// If the document contains meta data
	// the content will be between the description
	// and the meta data. If not the content
	// will go up to the end of the document.
	endLine = 0
	metaDataExists, metaDataLineNumber := doc.locateMetaData()
	if metaDataExists {
		endLine = metaDataLineNumber - 1
	} else {
		endLine = len(doc.rawLines)
	}

	// All lines between the start- and endLine are content
	return true, startLine, endLine
}
