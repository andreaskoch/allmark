package parser

import (
	"github.com/andreaskoch/docs/repository"
)

// Parse an item with a title, description and content
func parseDocumentLikeItem(item *repository.Item, lines []string) *repository.Item {

	// meta data
	item, lines = ParseMetaData(item, lines)

	// title
	item.Title, lines = getMatchingValue(lines, TitlePattern)

	// description
	item.Description, lines = getMatchingValue(lines, DescriptionPattern)

	// content
	item.Content = getContent(lines)

	return item
}
