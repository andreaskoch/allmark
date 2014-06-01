// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/content"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger/console"
	"github.com/andreaskoch/allmark2/common/logger/loglevel"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess/filesystem"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml"
	"github.com/andreaskoch/allmark2/services/initialization"
	"github.com/andreaskoch/allmark2/services/parsing"
	"github.com/andreaskoch/allmark2/services/search"
	"github.com/andreaskoch/allmark2/ui/web/server"
	"os"
	"path/filepath"
	"strings"
)

const (
	CommandNameInit  = "init"
	CommandNameServe = "serve"
)

func main() {

	serve := func(repositoryPath string) {

		// logger
		logger := console.New(loglevel.Debug)

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

		// repository (fallback) name
		repositoryName := filepath.Base(repositoryPath)

		// item index
		index := index.New(logger, repositoryName)

		// search
		itemSearch := search.NewItemSearch(logger, index)

		// read the repository
		for itemEvent := range repository.Items() {

			// validate event
			if itemEvent.Error != nil {
				logger.Warn("%s", itemEvent.Error)
				continue
			}

			repositoryItem := itemEvent.Item
			if repositoryItem == nil {
				logger.Warn("Repository item is empty.")
				continue
			}

			// define a process function
			processRepositoryItem := func() {
				// parse item
				if item, err := parser.Parse(repositoryItem); err == nil {
					// register the item at the index
					index.Add(item)
				} else {
					logger.Warn("%s", err.Error())
				}
			}

			// watch for changes
			go func() {
				for changeEvent := range repositoryItem.ChangeEvent() {

					switch changeEvent {

					case content.TypeChanged:
						{
							logger.Info("Item %q changed.", repositoryItem)
							processRepositoryItem()
						}

					case content.TypeMoved:
						{
							logger.Info("Item %q moved.", repositoryItem)
							break
						}

					}

				}
			}()

			processRepositoryItem()

		}

		// update the full-text search index
		itemSearch.Update()

		// server
		server, err := server.New(logger, config, index, converter, itemSearch)
		if err != nil {
			logger.Fatal("Unable to instantiate a server. Error: %s", err.Error())
		}

		// start the server
		if result := <-server.Start(); result != nil {
			logger.Info("%s", result)
		}
	}

	init := func(repositoryPath string) {
		if success, err := initialization.Initialize(repositoryPath); !success {
			fmt.Println(err)
		}
	}

	parseCommandLineArguments(os.Args, func(commandName, repositoryPath string) (commandWasFound bool) {
		switch strings.ToLower(commandName) {
		case CommandNameInit:
			init(repositoryPath)

		case CommandNameServe:
			serve(repositoryPath)

		default:
			return false
		}

		return true
	})
}

func parseCommandLineArguments(args []string, commandHandler func(commandName, repositoryPath string) (commandWasFound bool)) {

	// check if the mandatory amount of
	// command line parameters has been
	// supplied. If not, print usage information.
	if len(args) < 2 {
		printUsageInformation(args)
		return
	}

	// Read the repository path parameters
	var repositoryPath string
	if len(args) > 2 {

		// use supplied repository path
		repositoryPath = args[2]

		if isFile, _ := fsutil.IsFile(repositoryPath); isFile {
			repositoryPath = filepath.Dir(repositoryPath)
		}

	} else {

		// use the current directory
		repositoryPath = fsutil.GetWorkingDirectory()

	}

	// validate the supplied repository paths
	if !fsutil.PathExists(repositoryPath) {
		fmt.Fprintf(os.Stderr, "The specified repository paths %q is does not exist.", repositoryPath)
		return
	}

	// Read the command parameter and execute the command handler
	commandName := strings.ToLower(args[1])
	if commandWasFound := commandHandler(commandName, repositoryPath); !commandWasFound {
		printUsageInformation(args)
	}
}

// Print usage information
func printUsageInformation(args []string) {
	executeableName := args[0]

	fmt.Fprintf(os.Stderr, "%s - %s\n", executeableName, "A markdown web server and renderer")
	fmt.Fprintf(os.Stderr, "\nUsage:\n%s %s %s\n", executeableName, "<command>", "<repository path>")
	fmt.Fprintf(os.Stderr, "\nAvailable commands:\n")
	fmt.Fprintf(os.Stderr, "  %7s  %s\n", CommandNameInit, "Initialize the configuration")
	fmt.Fprintf(os.Stderr, "  %7s  %s\n", CommandNameServe, "Start serving the supplied repository via HTTP")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Fork me on GitHub %q\n", "https://github.com/andreaskoch/allmark")

	os.Exit(2)
}
