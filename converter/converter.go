package converter

import (
	"fmt"
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/repository"
)

func Convert(item *repository.Item) (*parser.ParsedItem, error) {

	if item.IsVirtual() {

		return convertVirtualItem(item)
	}

	return convertPhysicalItem(item)
}

func convertVirtualItem(item *repository.Item) (*parser.ParsedItem, error) {
	return parser.Parse(item)
}

func convertPhysicalItem(item *repository.Item) (*parser.ParsedItem, error) {

	// parse
	parsedItem, err := parser.Parse(item)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	// convert content
	parsedItem.ConvertedContent = toHtml(item, parsedItem.RawContent)

	return parsedItem, nil
}
