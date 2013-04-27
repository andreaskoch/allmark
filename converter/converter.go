package converter

import (
	"fmt"
	"github.com/andreaskoch/allmark/converter/html"
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"os"
)

func Convert(item *repository.Item, targetFormat string) (*parser.ParsedItem, error) {

	// open the file
	file, err := os.Open(item.Path())
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	defer file.Close()

	// get the raw lines
	lines := util.GetLines(file)

	// parse
	parsedItem, err := parser.Parse(lines, item)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	// convert content
	switch targetFormat {
	default:
		parsedItem.ConvertedContent = html.Convert(item, parsedItem.RawContent)
	}

	return parsedItem, nil
}
