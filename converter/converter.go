package converter

import (
	"fmt"
	"github.com/andreaskoch/allmark/converter/html"
	"github.com/andreaskoch/allmark/parser"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"os"
)

type Converter func() (*parser.Result, error)

func New(item *repository.Item, targetFormat string) Converter {

	// open the file
	file, err := os.Open(item.Path())
	if err != nil {
		return func() (*parser.Result, error) {
			return nil, fmt.Errorf("%s", err)
		}
	}

	defer file.Close()

	// get the raw lines
	lines := util.GetLines(file)

	// parse
	parsedItem, err := parser.Parse(lines, item.Type)
	if err != nil {
		return func() (*parser.Result, error) {
			return nil, fmt.Errorf("%s", err)
		}
	}

	// convert content
	parsedItem.ConvertedContent = html.Converter(item, parsedItem.RawContent)

	return func() (*parser.Result, error) {
		return parsedItem, nil
	}
}
