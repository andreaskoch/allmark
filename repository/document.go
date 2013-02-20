package repository

import (
	"andyk/docs/date"
	"andyk/docs/pattern"
	"log"
	"regexp"
	"strings"
)

type Document struct {
	Title       string
	Description string
	Content     string
	MetaData    MetaData
	Hash        string

	pattern  DocumentPattern
	rawLines []string
}

// CreateDocument returns a new Document from the given Item.
func CreateDocument(repositoryItem *Item) *Document {
	doc := Document{
		Hash:     repositoryItem.GetHash(),
		pattern:  NewDocumentPattern(),
		rawLines: repositoryItem.GetLines(),
	}

	// parse
	return doc.parse()
}

// getLastElement retrn the last element of a string array.
func getLastElement(array []string) string {
	if array == nil {
		return ""
	}

	return array[len(array)-1]
}

// parse starts the parsing of the current document.
// All regognized blocks will be assigned.
func (doc *Document) parse() *Document {
	return doc.setTitle()
}

// setTitle checks if the current Document
// contains a title and if yes, assigns it.
func (doc *Document) setTitle() *Document {
	titleLocation := doc.locateTitle()
	if !titleLocation.Found {
		return doc
	}

	// assemble title
	titleText := strings.TrimSpace(getLastElement(titleLocation.Matches))
	doc.Title = titleText

	return doc.setDescription()
}

// setDescription checks if the current Document
// contains a description block and if yes, assigns it.
func (doc *Document) setDescription() *Document {
	descriptionLocation := doc.locateDescription()
	if !descriptionLocation.Found {
		return doc
	}

	// assemble description
	descriptionText := strings.TrimSpace(getLastElement(descriptionLocation.Matches))
	doc.Description = descriptionText

	return doc.setContent()
}

// setContent checks if the current Document
// contains a content block and if yes, assigns it.
func (doc *Document) setContent() *Document {
	contentLocation := doc.locateContent()
	if !contentLocation.Found {
		return doc
	}

	// assemble content
	rawContent := strings.TrimSpace(strings.Join(contentLocation.Matches, "\n"))
	doc.Content = rawContent

	return doc.setMetaData()
}

// setMetaData checks if the current Document
// does contain meta data and if yes, assigns it.
func (doc *Document) setMetaData() *Document {
	metaDataLocation := doc.locateMetaData()
	if !metaDataLocation.Found {
		return doc
	}

	// assemble meta data
	var metaData MetaData

	for _, line := range metaDataLocation.Matches {
		isKeyValuePair, matches := pattern.IsMatch(line, doc.pattern.MetaData)

		// skip if line is not a key-value pair
		if !isKeyValuePair {
			continue
		}

		key := strings.ToLower(strings.TrimSpace(matches[1]))
		value := strings.TrimSpace(matches[2])

		switch strings.ToLower(key) {

		case "language":
			{
				metaData.Language = value
				break
			}

		case "date":
			{
				date, err := date.ParseIso8601Date(value)
				if err == nil {
					metaData.Date = date
				}
				break
			}

		case "tags":
			{
				metaData.Tags = GetTagsFromValue(value)
				break
			}

		}
	}

	// assign meta data
	doc.MetaData = metaData

	return doc
}

func GetTagsFromValue(value string) []string {
	rawTags := strings.Split(value, ",")
	tags := make([]string, 0, 1)

	for _, tag := range rawTags {
		trimmedTag := strings.TrimSpace(tag)
		if trimmedTag != "" {
			tags = append(tags, trimmedTag)
		}

	}

	return tags
}

// DocumentPattern contains a set of regular expression
// for parsing documents.
type DocumentPattern struct {
	EmptyLine      regexp.Regexp
	Title          regexp.Regexp
	Description    regexp.Regexp
	HorizontalRule regexp.Regexp
	MetaData       regexp.Regexp
}

// NewDocumentPattern returns a new DocumentPattern
// for parsing documents.
func NewDocumentPattern() DocumentPattern {
	// Lines which contain nothing but white space characters
	// or no characters at all.
	emptyLineRegexp := regexp.MustCompile("^\\s*$")

	// Lines which a start with a hash, followed by zero or more
	// white space characters, followed by text.
	titleRegexp := regexp.MustCompile("\\s*#\\s*(\\w.+)")

	// Lines which start with text
	descriptionRegexp := regexp.MustCompile("^\\w.+")

	// Lines which nothing but dashes
	horizontalRuleRegexp := regexp.MustCompile("^-{2,}")

	// Lines with a "key: value" syntax
	metaDataRegexp := regexp.MustCompile("^(\\w+):\\s*(\\w.+)$")

	return DocumentPattern{
		EmptyLine:      *emptyLineRegexp,
		Title:          *titleRegexp,
		Description:    *descriptionRegexp,
		HorizontalRule: *horizontalRuleRegexp,
		MetaData:       *metaDataRegexp,
	}
}

// LineRange contains a Start- and a End line number.
type LineRange struct {
	Start int
	End   int
}

// NewLineRange returns a new LineRange
// with the given start and end.
func NewLineRange(start int, end int) LineRange {
	if start < 0 || end < 0 || (start > end) {
		log.Panicf("Invalid start and end values for a LineRange. Start: %v, End: %v", start, end)
	}

	return LineRange{
		Start: start,
		End:   end,
	}
}

// A MatchResult represents the result of a pattern matching
// process on the content of an document.
// It indicates whether the pattern was found and if yet,
// the lines in which it was located and the matched text.
type MatchResult struct {
	Found   bool
	Lines   LineRange
	Matches []string
}

// Found create a new MatchResult which represents
// a successful match.
func Found(firstLine int, lastLine int, matches []string) *MatchResult {
	return &MatchResult{
		Found:   true,
		Lines:   NewLineRange(firstLine, lastLine),
		Matches: matches,
	}
}

// NotFound create a new MatchResult which represents
// an unsuccessful match.
func NotFound() *MatchResult {
	return &MatchResult{
		Found: false,
		Lines: NewLineRange(-1, -1),
	}
}

// locateTitle checks if the current Document
// contains a title.
func (doc *Document) locateTitle() *MatchResult {

	// In order to be the "title" the line must either
	// be empty or match the title pattern.

	for lineNumber, line := range doc.rawLines {

		lineMatchesTitlePattern, matches := pattern.IsMatch(line, doc.pattern.Title)
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

// locateDescription checks if the current Document
// contains a description.
func (doc *Document) locateDescription() *MatchResult {

	// The description must be preceeded by a title
	title := doc.locateTitle()
	if !title.Found {
		return NotFound()
	}

	// If the document has no more lines than the line
	// in which the title has been located, there
	// will be no room for a description
	startLine := title.Lines.End + 1
	if len(doc.rawLines) <= startLine {
		return NotFound()
	}

	// In order to be a "description" the line must either
	// be empty or match the description pattern.
	for relativeLineNumber, line := range doc.rawLines[startLine:] {

		lineMatchesDescriptionPattern, matches := pattern.IsMatch(line, doc.pattern.Description)
		if lineMatchesDescriptionPattern {
			absoluteLineNumber := startLine + relativeLineNumber
			return Found(absoluteLineNumber, absoluteLineNumber, matches)
		}

		lineIsEmpty := doc.pattern.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			break
		}
	}

	return NotFound()
}

// locateMetaData checks if the current Document
// contains meta data.
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
	for _, line := range doc.rawLines[metaDataStartLine:] {

		lineMatchesMetaDataPattern := doc.pattern.MetaData.MatchString(line)
		if lineMatchesMetaDataPattern {

			endLine := len(doc.rawLines)
			return Found(metaDataStartLine, endLine, doc.rawLines[metaDataStartLine:endLine])

		}

		lineIsEmpty := doc.pattern.EmptyLine.MatchString(line)
		if !lineIsEmpty {
			return NotFound()
		}

	}

	return NotFound()
}

// locateContent checks if the current Document
// contains content.
func (doc *Document) locateContent() *MatchResult {

	// Content must be preceeded by a description
	description := doc.locateDescription()
	if !description.Found {
		return NotFound()
	}

	// If the document has no more lines than the line
	// in which the description has been located, there
	// will be no room for content
	startLine := description.Lines.End + 1
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
