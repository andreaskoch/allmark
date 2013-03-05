package parser

import (
	"github.com/andreaskoch/docs/indexer"
)

func ParseRepository(item *indexer.Item, lines []string) *indexer.Item {

	// meta data
	item, lines = ParseMetaData(item, lines)

	// title
	title, lines := getMatchingValue(lines, TitlePattern)
	item.AddBlock("title", title)

	// description
	description, lines := getMatchingValue(lines, DescriptionPattern)
	item.AddBlock("description", description)

	// content
	item.AddBlock("content", getContent(lines))

	return item
}
