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

type LineSet struct {
	Start int
	End   int
}

func NewLineSet(start int, end int) LineSet {
	return LineSet{
		Start: start,
		End:   end,
	}
}

type MatchResult struct {
	Found   bool
	Lines   LineSet
	Matches []string
}

func Found(firstLine int, lastLine int, matches []string) *MatchResult {
	return &MatchResult{
		Found: true,
		Lines: NewLineSet(firstLine, lastLine),
	}
}

func NotFound() *MatchResult {
	return &MatchResult{
		Found: false,
		Lines: NewLineSet(-1, -1),
	}
}

func IsMatch(line string, pattern regexp.Regexp) (isMatch bool, matches []string) {
	matches = pattern.FindStringSubmatch(line)
	return matches != nil, matches
}

// Check if the current Document contains a title
func (doc *Document) locateTitle() *MatchResult {

	// In order to be the "title" the line must either
	// be empty or match the title pattern.

	for lineNumber, line := range doc.rawLines {

		lineMatchesTitlePattern, matches := IsMatch(line, doc.pattern.Title)
		if lineMatchesTitlePattern {
			return Found(lineNumber, lineNumber, matches)
		}

		lineIsEmpty := doc.pattern.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return NotFound()
}

// Check if the current Document contains a description
func (doc *Document) locateDescription() *MatchResult {

	// The description must be preceeded by a title
	title := doc.locateTitle()
	if !title.Found {
		return NotFound()
	}

	// If the document has no more lines than the line
	// in which the title has been located, there
	// will be no room for a description
	startLine := title.Lines.Start + 1
	if len(doc.rawLines) <= startLine {
		return NotFound()
	}

	// In order to be a "description" the line must either
	// be empty or match the description pattern.
	for lineNumber, line := range doc.rawLines[startLine:] {

		lineMatchesDescriptionPattern, matches := IsMatch(line, doc.pattern.Description)
		if lineMatchesDescriptionPattern {
			return Found(lineNumber, lineNumber, matches)
		}

		lineIsEmpty := doc.pattern.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return NotFound()
}

// Check if the current Document contains meta data
func (doc *Document) locateMetaData() *MatchResult {

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
		return NotFound()
	}

	// If the document has no more lines than
	// the last found horizontal rule there is no
	// room for meta data
	metaDataStartLine := lastFoundHorizontalRulePosition + 1
	if len(doc.rawLines) <= metaDataStartLine {
		return NotFound()
	}

	// Check if the last horizontal rule is followed
	// either by white space or be meta data
	for lineNumber, line := range doc.rawLines[metaDataStartLine:] {

		lineMatchesMetaDataPattern := doc.pattern.MetaData.MatchString(line)
		if lineMatchesMetaDataPattern {
			return Found(metaDataStartLine, len(doc.rawLines), doc.rawLines[metaDataStartLine:])
		}

		lineIsEmpty := doc.pattern.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			return NotFound()
		}

	}

	return NotFound()
}

// Check if the current Document contains content
func (doc *Document) locateContent() *MatchResult {

	// Content must be preceeded by a description
	description := doc.locateDescription()
	if !description.Found {
		return NotFound()
	}

	// If the document has no more lines than the line
	// in which the description has been located, there
	// will be no room for content
	startLine := description.Lines.Start + 1
	if len(doc.rawLines) <= startLine {
		return NotFound()
	}

	// If the document contains meta data
	// the content will be between the description
	// and the meta data. If not the content
	// will go up to the end of the document.
	endLine := 0
	metaData := doc.locateMetaData()
	if metaData.Found {
		endLine = metaData.Lines.Start - 1
	} else {
		endLine = len(doc.rawLines)
	}

	// All lines between the start- and endLine are content
	return Found(startLine, endLine, doc.rawLines[startLine:endLine])
}
