package parser

import (
	"github.com/andreaskoch/docs/indexer"
)

func ParseDocument(item *indexer.Item, lines []string) *indexer.Item {

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
