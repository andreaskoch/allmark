// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/parser/cleanup"
	"github.com/andreaskoch/allmark2/services/parser/document"
	"github.com/andreaskoch/allmark2/services/parser/message"
	"github.com/andreaskoch/allmark2/services/parser/presentation"
	"github.com/andreaskoch/allmark2/services/parser/typedetection"
	"io"
)

type Parser struct {
	logger logger.Logger
}

func New(logger logger.Logger) (Parser, error) {
	return Parser{
		logger: logger,
	}, nil
}

func (parser *Parser) ParseItem(item *dataaccess.Item) (*model.Item, error) {

	if item == nil {
		return nil, fmt.Errorf("Cannot parse an empty item.")
	}

	parser.logger.Debug("Parsing item %q", item.String())
	route := item.Route()

	// convert the files
	files := parser.convertFiles(item.Files())

	// create a new item model
	itemModel, err := model.NewItem(route, files)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert Item %q. Error: %s", item, err)
	}

	// capture the last modified date
	lastModifiedDate, err := item.LastModified()

	// fetch the item data
	byteBuffer := new(bytes.Buffer)
	dataWriter := bufio.NewWriter(byteBuffer)
	contentReader := func(content io.ReadSeeker) error {
		_, err := io.Copy(dataWriter, content)
		dataWriter.Flush()
		return err
	}

	if err := item.Data(contentReader); err != nil {
		return nil, err
	}

	data := byteBuffer.Bytes()

	lines := getLines(bytes.NewReader(data))

	// cleanup the markdown before parser it
	lines = cleanup.Cleanup(lines)

	// detect the item type
	switch itemModel.Type = typedetection.DetectType(lines); itemModel.Type {

	case model.TypeDocument, model.TypeLocation, model.TypeRepository:
		{
			if _, err := document.Parse(itemModel, lastModifiedDate, lines); err != nil {
				return nil, fmt.Errorf("Unable to parse item %q (Type: %s, Error: %s)", item, itemModel.Type, err.Error())
			}
		}

	case model.TypePresentation:
		{
			if err := presentation.Parse(itemModel, lastModifiedDate, lines); err != nil {
				return nil, fmt.Errorf("Unable to parse item %q (Type: %s, Error: %s)", item, itemModel.Type, err.Error())
			}
		}

	case model.TypeMessage:
		{
			if err := message.Parse(itemModel, lastModifiedDate, lines); err != nil {
				return nil, fmt.Errorf("Unable to parse item %q (Type: %s, Error: %s)", item, itemModel.Type, err.Error())
			}
		}

	default:
		return nil, fmt.Errorf("Cannot parse item %q. Unknown item type.", item)

	}

	return itemModel, nil
}

func (parser *Parser) ParseFile(file *dataaccess.File) (*model.File, error) {

	convertedFile, err := model.NewFromDataAccess(file)
	if err != nil {
		return nil, fmt.Errorf("Unable to convert file %q. Error: %s", file, err.Error())
	}

	return convertedFile, nil
}

func (parser *Parser) convertFiles(dataaccessFiles []*dataaccess.File) []*model.File {

	convertedFiles := make([]*model.File, 0, len(dataaccessFiles))

	for _, file := range dataaccessFiles {

		if convertedFile, err := parser.ParseFile(file); err == nil {
			convertedFiles = append(convertedFiles, convertedFile)
		}

	}

	return convertedFiles
}
