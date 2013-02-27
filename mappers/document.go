package mappers

import (
	"andyk/docs/parser"
	"andyk/docs/viewmodel"
)

func GetDocument(parsedItem parser.ParsedItem) viewmodel.Document {
	return viewmodel.Document{
		Title:       parsedItem.GetElementValue("title"),
		Description: parsedItem.GetElementValue("description"),
		Content:     parsedItem.GetElementValue("content"),
		LanguageTag: getTwoLetterLanguageCode(parsedItem.MetaData.Language),
	}
}
