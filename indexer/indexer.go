package indexer

import (
	"andyk/docs/model"
)

func Index() model.Document {

	var doc model.Document
	doc.Path = "Test"
	doc.Title = "Test"
	doc.Description = "Description"
	doc.Content = "Content"
	doc.Language = "en-US"

	return doc
}
