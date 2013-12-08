// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger/console"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess/filesystem"
	"github.com/andreaskoch/allmark2/services/converter/markdowntohtml"
	"github.com/andreaskoch/allmark2/services/parser"
)

func main() {

	// logger
	logger := console.New()

	// data access
	repository, err := filesystem.NewRepository(logger, fsutil.GetWorkingDirectory())
	if err != nil {
		panic(err)
	}

	// parser
	parser, err := parser.New(logger)
	if err != nil {
		panic(err)
	}

	// converter
	converter, err := markdowntohtml.New(logger)
	if err != nil {
		panic(err)
	}

	// read the repository
	itemEvents := repository.GetItems()
	for itemEvent := range itemEvents {

		if itemEvent.Error != nil {
			logger.Warn("%s", itemEvent.Error)
		}

		if itemEvent.Item == nil {
			continue
		}

		// parse item
		item, err := parser.Parse(itemEvent.Item)
		if err != nil {
			logger.Warn("Unable to parse item %q. Error: %s", itemEvent.Item, err)
			continue
		}

		// convert item
		converter.Convert(item)
		fmt.Println(item.Content)
	}
}
