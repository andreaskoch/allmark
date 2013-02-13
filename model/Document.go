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
	return doc.setTitle().setDescription()
}

func (doc *Document) setTitle() *Document {
	titleRegexp := regexp.MustCompile("\\s*#\\s*(.+)")

	for lineNumber, line := range doc.rawLines[doc.lastFindIndex:] {
		matches := titleRegexp.FindStringSubmatch(line)

		if len(matches) == 2 {
			doc.lastFindIndex = lineNumber
			doc.Title = matches[1]
			return doc
		}
	}

	return doc
}

func (doc *Document) setDescription() *Document {
	descriptionRegexp := regexp.MustCompile("^\\w.+")

	for lineNumber, line := range doc.rawLines[doc.lastFindIndex:] {
		matches := descriptionRegexp.FindStringSubmatch(line)

		if len(matches) == 1 {
			doc.lastFindIndex = lineNumber
			doc.Description = matches[0]
			return doc
		}
	}

	return doc
}
