// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/andreaskoch/allmark/common/logger"
	"github.com/andreaskoch/allmark/dataaccess"
	"github.com/andreaskoch/allmark/model"
	"github.com/andreaskoch/allmark/services/parser/cleanup"
	"github.com/andreaskoch/allmark/services/parser/document"
	"github.com/andreaskoch/allmark/services/parser/presentation"
	"github.com/andreaskoch/allmark/services/parser/typedetection"
)

type Parser struct {
	logger logger.Logger
}

func New(logger logger.Logger) (Parser, error) {
	return Parser{
		logger: logger,
	}, nil
}

func (parser *Parser) ParseItem(item dataaccess.Item) (*model.Item, error) {

	if item == nil {
		return nil, fmt.Errorf("Cannot parse an empty item.")
	}

	parser.logger.Debug("Parsing item %q", item.String())
	route := item.Route()

	// convert the files
	files := parser.convertFiles(item.Files())

	// create a new item model
	itemModel := model.NewItem(route, files, item.Type())

	// capture the last modified date
	lastModifiedDate, err := item.LastModified()
	if err != nil {
		return nil, fmt.Errorf("Cannot determine last modified date for item %q. Error: %s", item, err.Error())
	}

	// fetch the item data
	data, err := getItemData(item)
	if err != nil {
		return nil, fmt.Errorf("Cannot get data from item %q. Error: %s", item, err.Error())
	}

	// capture the markdown
	itemModel.Markdown = string(data)

	// split the markdown content into separate lines
	lines := getLines(bytes.NewReader(data))
	lines = cleanup.Cleanup(lines)

	// detect the item type
	switch itemModel.Type = typedetection.DetectType(lines); itemModel.Type {

	case model.TypeDocument, model.TypeRepository:
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

	default:
		return nil, fmt.Errorf("Cannot parse item %q. Unknown item type.", item)

	}

	// item hash
	hash, err := item.Hash()
	if err != nil {
		return nil, fmt.Errorf("Unable to determine the hash for item %q. Error: %s", item, err.Error())
	}

	itemModel.Hash = hash

	return itemModel, nil
}

func (parser *Parser) ParseFile(file dataaccess.File) (*model.File, error) {

	return &model.File{
		file,
	}, nil

}

func (parser *Parser) convertFiles(dataaccessFiles []dataaccess.File) []*model.File {

	convertedFiles := make([]*model.File, 0, len(dataaccessFiles))

	for _, file := range dataaccessFiles {

		if convertedFile, err := parser.ParseFile(file); err == nil {
			convertedFiles = append(convertedFiles, convertedFile)
		}

	}

	return convertedFiles
}

func getItemData(item dataaccess.Item) ([]byte, error) {

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

	return byteBuffer.Bytes(), nil
}
