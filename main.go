// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/andreaskoch/allmark2/common/logger/console"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess/filesystem"
	"github.com/andreaskoch/allmark2/services/converter/markdowntohtml"
	"github.com/andreaskoch/allmark2/services/parser"
	"github.com/andreaskoch/allmark2/ui/web/server"
)

func main() {

	// logger
	logger := console.New()

	// data access
	repository, err := filesystem.NewRepository(logger, fsutil.GetWorkingDirectory())
	if err != nil {
		logger.Fatal("Unable to create a repository. Error: %s", err)
	}

	// parser
	parser, err := parser.New(logger)
	if err != nil {
		logger.Fatal("Unable to instantiate a parser. Error: %s", err)
	}

	// converter
	converter, err := markdowntohtml.New(logger)
	if err != nil {
		logger.Fatal("Unable to instantiate a converter. Error: %s", err)
	}

	// server
	server, err := server.New(logger)
	if err != nil {
		logger.Fatal("Unable to instantiate a server. Error: %s", err)
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

		// send item to server
		server.Serve(item)
	}
}
