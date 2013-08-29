// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"
	"github.com/andreaskoch/allmark/parser/document"
	"github.com/andreaskoch/allmark/parser/message"
	"github.com/andreaskoch/allmark/parser/metadata"
	"github.com/andreaskoch/allmark/parser/presentation"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/types"
	"github.com/andreaskoch/allmark/util"
	"os"
	"path/filepath"
	"strings"
)

func Parse(item *repository.Item) (*repository.Item, error) {
	if item.IsVirtual() {
		return parseVirtual(item)
	}

	return parsePhysical(item)
}

func getFallbackTitle(item *repository.Item) string {
	return filepath.Base(item.Directory())
}

func parseVirtual(item *repository.Item) (*repository.Item, error) {

	if item == nil {
		return nil, fmt.Errorf("Cannot create meta data from nil.")
	}

	// get the item title
	title := getFallbackTitle(item)

	// create the meta data
	metaData, err := metadata.New(item)
	if err != nil {
		return nil, err
	}

	item.Title = title
	item.MetaData = metaData

	return item, nil
}

func parsePhysical(item *repository.Item) (*repository.Item, error) {

	// open the file
	file, err := os.Open(item.Path())
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	defer file.Close()

	// get the raw lines
	lines := util.GetLines(file)

	// determine the fallback title
	fallbackTitle := getFallbackTitle(item)

	// parse the meta data
	fallbackItemTypeFunc := func() string {
		return getItemTypeFromFilename(item.Path())
	}

	item.MetaData, lines = metadata.Parse(item, lines, fallbackItemTypeFunc)

	// parse the content
	switch itemType := item.MetaData.ItemType; itemType {

	case types.RepositoryItemType, types.DocumentItemType:
		{
			if success, err := document.Parse(item, lines, fallbackTitle); success {
				return item, nil
			} else {
				return nil, err
			}
		}

	case types.PresentationItemType:
		{
			if success, err := presentation.Parse(item, lines, fallbackTitle); success {
				return item, nil
			} else {
				return nil, err
			}
		}

	case types.MessageItemType:
		{
			if success, err := message.Parse(item, lines, fallbackTitle); success {
				return item, nil
			} else {
				return nil, err
			}
		}

	default:
		return nil, fmt.Errorf("Item %q (type: %s) cannot be parsed.", item.Path(), itemType)

	}

	panic("Unreachable")
}

func getItemTypeFromFilename(filenameOrPath string) string {

	extension := filepath.Ext(filenameOrPath)
	filenameWithExtension := filepath.Base(filenameOrPath)
	filename := filenameWithExtension[0:(strings.LastIndex(filenameWithExtension, extension))]

	switch strings.ToLower(filename) {
	case types.DocumentItemType:
		return types.DocumentItemType

	case types.PresentationItemType:
		return types.PresentationItemType

	case types.MessageItemType:
		return types.MessageItemType

	case types.RepositoryItemType:
		return types.RepositoryItemType

	default:
		return types.DocumentItemType
	}

	return types.UnknownItemType
}
