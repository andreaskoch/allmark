package model

import (
	"regexp"
)

type Document struct {
	Title       string
	Description string
	Content     string
	Hash        string

	rawLines []string
}

func CreateDocument(repositoryItem *RepositoryItem) *Document {
	doc := Document{
		Hash:     repositoryItem.GetHash(),
		rawLines: repositoryItem.GetLines(),
	}

	// parse
	return doc.parse()
}

func (doc *Document) parse() *Document {
	return setTitle(doc)
}

func setTitle(doc *Document) *Document {
	titleRegexp := regexp.MustCompile("\\s*#\\s*(.+)")

	for lineNumber, line := range doc.rawLines {
		matches := titleRegexp.FindStringSubmatch(line)

		// line must match title pattern
		lineMatchesTitlePattern := len(matches) == 2
		if lineMatchesTitlePattern {

			// is first line or all previous lines are empty
			if lineNumber == 0 || linesMeetCondition(doc.rawLines[0:lineNumber], regexp.MustCompile("^\\s*$")) {

				doc.Title = matches[1]
				return setDescription(doc, lineNumber+1)

			}
		}
	}

	return doc
}

func setDescription(doc *Document, startLine int) *Document {
	if startLine > len(doc.rawLines) {
		return doc
	}

	descriptionRegexp := regexp.MustCompile("^\\w.+")

	for lineNumber, line := range doc.rawLines[startLine:] {
		matches := descriptionRegexp.FindStringSubmatch(line)

		// line must match description pattern
		lineMatchesDescriptionPattern := len(matches) == 1
		if lineMatchesDescriptionPattern {
			doc.Description = matches[0]
			return setContent(doc, lineNumber+1)
		}
	}

	return doc
}

func setContent(doc *Document, startLine int) *Document {
	if startLine > len(doc.rawLines) {
		return doc
	}

	content := ""

	for _, line := range doc.rawLines[startLine:] {
		content += line + "\n"
	}

	doc.Content = content

	return doc
}

func linesMeetCondition(lines []string, condition *regexp.Regexp) bool {

	for _, line := range lines {
		if !condition.MatchString(line) {
			return false
		}
	}

	return true
}
