package converter

import (
	"fmt"
	"github.com/andreaskoch/allmark/converter/html"
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/repository"
)

func Convert(item *repository.Item) (*repository.Item, error) {

	if item.IsVirtual() {

		return convertVirtualItem(item)
	}

	return convertPhysicalItem(item)
}

func convertVirtualItem(item *repository.Item) (*repository.Item, error) {
	return parser.Parse(item)
}

func convertPhysicalItem(item *repository.Item) (*repository.Item, error) {

	// parse
	_, err := parser.Parse(item)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	// convert content
	html.ToHtml(item)

	return item, nil
}
