// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/parser/document"
	"github.com/andreaskoch/allmark2/services/parser/typedetection"
)

type Parser struct {
	logger logger.Logger
}

func New(logger logger.Logger) (*Parser, error) {
	return &Parser{
		logger: logger,
	}, nil
}

func (parser *Parser) Parse(item *dataaccess.Item) (*model.Item, error) {

	route := item.Route()

	// convert the files
	files := parser.convertFiles(item.Files())

	itemModel, err := model.NewItem(route, files)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert Item %q. Error: %s", item, err)
	}

	// fetch the item content
	getContent := item.ContentProvider()
	content, _ := getContent()
	lines := getLines(bytes.NewReader(content))

	// detect the item type
	switch itemType := typedetection.DetectType(lines); itemType {

	case model.Document:
		{
			if err := document.Parse(itemModel, lines); err != nil {
				return nil, fmt.Errorf("Unable to parse document title of Item %q. Error: %s", item, err)
			}
		}

	case model.Unknown:
		return nil, fmt.Errorf("Unknown item type for Item %q.", item)

	default:
		panic("Unreachable")

	}

	return itemModel, nil
}

func (parser *Parser) convertFiles(dataaccessFiles []*dataaccess.File) []*model.File {

	files := make([]*model.File, 0, len(dataaccessFiles))

	for _, file := range dataaccessFiles {

		fileModel, err := model.NewFile(file)
		if err != nil {
			parser.logger.Warn("Unable to convert file %q.", file)
			continue
		}

		files = append(files, fileModel)
	}

	return files
}
