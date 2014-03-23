// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger/console"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess/filesystem"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml"
	"github.com/andreaskoch/allmark2/services/parsing"
	"github.com/andreaskoch/allmark2/ui/web/server"
)

func main() {

	repositoryPath := fsutil.GetWorkingDirectory()

	// logger
	logger := console.New()

	// config
	config := config.Get(repositoryPath)

	// data access
	repository, err := filesystem.NewRepository(logger, repositoryPath)
	if err != nil {
		logger.Fatal("Unable to create a repository. Error: %s", err)
	}

	// parser
	parser, err := parsing.New(logger)
	if err != nil {
		logger.Fatal("Unable to instantiate a parser. Error: %s", err)
	}

	// converter
	converter, err := markdowntohtml.New(logger)
	if err != nil {
		logger.Fatal("Unable to instantiate a converter. Error: %s", err)
	}

	// item index
	itemIndex := index.CreateItemIndex(logger)

	// file index
	fileIndex := index.CreateFileIndex(logger)

	// server
	server, err := server.New(logger, config, converter, itemIndex, fileIndex)
	if err != nil {
		logger.Fatal("Unable to instantiate a server. Error: %s", err)
	}

	// serve theme files
	baseFolder := config.MetaDataFolder()
	themeFolder := config.ThemeFolder()
	fileIndex.AddFolder(baseFolder, themeFolder)

	// read the repository
	itemEvents := repository.GetItems()
	for itemEvent := range itemEvents {

		// validate event
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

		// register the item at the index
		itemIndex.Add(item)
	}

	// start the server
	if result := <-server.Start(); result != nil {
		logger.Info("%s", result)
	}
}
