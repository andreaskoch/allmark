package converter

import (
	"fmt"
	"github.com/andreaskoch/allmark/converter/html"
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/repository"
)

func Convert(item *repository.Item, targetFormat string) (*parser.ParsedItem, error) {

	if item.IsVirtual() {

		return convertVirtualItem(item, targetFormat)
	}

	return convertPhysicalItem(item, targetFormat)
}

func convertVirtualItem(item *repository.Item, targetFormat string) (*parser.ParsedItem, error) {
	return parser.Parse(item)
}

func convertPhysicalItem(item *repository.Item, targetFormat string) (*parser.ParsedItem, error) {

	// parse
	parsedItem, err := parser.Parse(item)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	// convert content
	switch targetFormat {
	default:
		parsedItem.ConvertedContent = html.ToHtml(item, parsedItem.RawContent)
	}

	return parsedItem, nil
}
