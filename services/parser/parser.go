// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/model"
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

	// convert the files
	files := parser.convertFiles(item.Files())

	itemModel, err := model.NewItem(item.Route(), files)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert item %q. Error: %s", item, err)
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
