package model

import (
	"regexp"
)

type Document struct {
	Title       string
	Description string
	Hash        string

	lastFindIndex int
	rawLines      []string
}

func CreateDocument(repositoryItem *RepositoryItem) *Document {
	doc := Document{
		Hash:          repositoryItem.GetHash(),
		rawLines:      repositoryItem.GetLines(),
		lastFindIndex: 0,
	}

	// parse
	return doc.parse()
}

func (doc *Document) parse() *Document {
	return doc.setTitle()
}

func (doc *Document) setTitle() *Document {
	titleRegexp := regexp.MustCompile("\\s*#\\s*(.+)")

	for lineNumber, line := range doc.rawLines[doc.lastFindIndex:] {
		matches := titleRegexp.FindStringSubmatch(line)

		// line must match title pattern
		lineMatchesTitlePattern := len(matches) == 2
		if lineMatchesTitlePattern {

			// is first line or all previous lines are empty
			if lineNumber == 0 || linesMeetCondition(doc.rawLines[0:lineNumber], regexp.MustCompile("^\\s*$")) {

				doc.lastFindIndex = lineNumber
				doc.Title = matches[1]
				return doc.setDescription()

			}
		}
	}

	return doc
}

func (doc *Document) setDescription() *Document {
	descriptionRegexp := regexp.MustCompile("^\\w.+")

	for lineNumber, line := range doc.rawLines[doc.lastFindIndex:] {
		matches := descriptionRegexp.FindStringSubmatch(line)

		// line must match description pattern
		lineMatchesDescriptionPattern := len(matches) == 1
		if lineMatchesDescriptionPattern {
			doc.lastFindIndex = lineNumber
			doc.Description = matches[0]
			return doc
		}
	}

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
