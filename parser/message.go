package parser

import (
	"github.com/andreaskoch/docs/indexer"
)

func parseMessage(item *indexer.Item, lines []string) *indexer.Item {

	// meta data
	item, lines = ParseMetaData(item, lines)

	// content
	item.Content = getContent(lines)

	return item
}
