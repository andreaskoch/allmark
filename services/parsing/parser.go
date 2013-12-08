// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parsing

import (
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/parsing/document"
	"github.com/andreaskoch/allmark2/services/parsing/message"
	"github.com/andreaskoch/allmark2/services/parsing/presentation"
	"github.com/andreaskoch/allmark2/services/parsing/typedetection"
	"time"
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

	parser.logger.Debug("Parsing item %q", item)

	route := item.Route()

	// convert the files
	files := parser.convertFiles(item.Files())

	// create a new item model
	itemModel, err := model.NewItem(route, files)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert Item %q. Error: %s", item, err)
	}

	// content provider
	contentProvider := item.ContentProvider()

	// capture the last modified date
	lastModifiedDate, err := contentProvider.LastModified()
	if err != nil {
		lastModifiedDate = time.Now() // fallback to the current time
	}

	// fetch the item data
	data, _ := contentProvider.Data()

	// todo: cleanup data
	// - replace \r\n with \n

	lines := getLines(bytes.NewReader(data))

	// detect the item type
	switch itemType := typedetection.DetectType(lines); itemType {

	case model.TypeDocument, model.TypeLocation, model.TypeRepository:
		{
			if err := document.Parse(itemModel, lastModifiedDate, lines); err != nil {
				return nil, fmt.Errorf("Unable to parse item %q. Error: %s", item, err)
			}
		}

	case model.TypePresentation:
		{
			if err := presentation.Parse(itemModel, lastModifiedDate, lines); err != nil {
				return nil, fmt.Errorf("Unable to parse item %q. Error: %s", item, err)
			}
		}

	case model.TypeMessage:
		{
			if err := message.Parse(itemModel, lastModifiedDate, lines); err != nil {
				return nil, fmt.Errorf("Unable to parse item %q. Error: %s", item, err)
			}
		}

	default:
		return nil, fmt.Errorf("Cannot parse item %q. Unknown item type.", item)

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
