package parser

import (
	"github.com/andreaskoch/docs/repository"
)

func parseMessage(item *repository.Item, lines []string) *repository.Item {

	// meta data
	item, lines = ParseMetaData(item, lines)

	// content
	item.Content = getContent(lines)

	return item
}
