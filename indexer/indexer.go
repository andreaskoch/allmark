package indexer

import (
	"andyk/docs/model"
	"time"
)

func Index() model.Document {

	var doc model.Document
	doc.Path = "Test"
	doc.Title = "Test"
	doc.Description = "Description"
	doc.Content = "Content"
	doc.Language = "en-US"
	doc.Date = time.Date(2013, 1, 13, 0, 0, 0, 0, time.UTC)

	return doc
}
