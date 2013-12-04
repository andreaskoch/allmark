// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger/console"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess/filesystem"
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

	// read the repository
	itemEvents, done := repository.GetItems()

	allItemsRetrieved := false
	for !allItemsRetrieved {
		select {
		case isDone := <-done:
			if isDone {
				allItemsRetrieved = true
			}

		case itemEvent := <-itemEvents:

			if itemEvent.Item != nil {

				// parse item
				logger.Info("Parsing item %q", itemEvent.Item)
				item, err := parser.Parse(itemEvent.Item)
				if err != nil {
					logger.Warn("Unable to parse item %q. Error: %s", itemEvent.Item, err)
					continue
				}

				logger.Info("Parsed item %q.", item.Title)
				fmt.Println(item.Description)
				fmt.Println("---------")
				fmt.Println(item.Content)
				fmt.Println("---------")
				fmt.Printf("%#v\n", item.MetaData)

			}

		}
	}
}
