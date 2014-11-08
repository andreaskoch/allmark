// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/logger/console"
	"github.com/andreaskoch/allmark2/common/logger/loglevel"
	"github.com/andreaskoch/allmark2/common/shutdown"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/dataaccess/filesystem"
	"github.com/andreaskoch/allmark2/services/converter/markdowntohtml"
	"github.com/andreaskoch/allmark2/services/initialization"
	"github.com/andreaskoch/allmark2/services/parser"
	"github.com/andreaskoch/allmark2/services/thumbnail"
	"github.com/andreaskoch/allmark2/web/server"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	CommandNameInit  = "init"
	CommandNameServe = "serve"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Handle CTRL-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case _ = <-c:
			{
				fmt.Println("Stopping")

				// Execute shutdown handlers
				shutdown.Shutdown()

				os.Exit(0)
			}
		}
	}()

	parseCommandLineArguments(os.Args, func(commandName, repositoryPath string) (commandWasFound bool) {
		switch strings.ToLower(commandName) {
		case CommandNameInit:
			initialize(repositoryPath)
			return true

		case CommandNameServe:
			serve(repositoryPath)
			return true

		default:
			return false
		}

		panic("Unreachable")
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

func serve(repositoryPath string) bool {

	config := *config.Get(repositoryPath)
	logger := console.New(loglevel.FromString(config.LogLevel))

	// data access
	repository, err := filesystem.NewRepository(logger, repositoryPath, config.Indexing.IntervalInSeconds)
	if err != nil {
		logger.Error("Unable to create a repository. Error: %s", err)
	}

	// thumbnail conversion service
	thumbnail.NewConversionService(logger, config, repository)

	// parser
	itemParser, err := parser.New(logger)
	if err != nil {
		logger.Error("Unable to instantiate a parser. Error: %s", err)
	}

	// converter
	converter, err := markdowntohtml.New(logger)
	if err != nil {
		logger.Error("Unable to instantiate a converter. Error: %s", err)
	}

	// server
	server, err := server.New(logger, config, repository, itemParser, converter)
	if err != nil {
		logger.Error("Unable to instantiate a server. Error: %s", err.Error())
		return false
	}

	if result := <-server.Start(); result != nil {
		logger.Error("%s", result)
		return false
	}

	return true
}

func initialize(repositoryPath string) bool {

	config := config.Get(repositoryPath)
	logger := console.New(loglevel.FromString(config.LogLevel))

	if success, err := initialization.Initialize(repositoryPath); !success {
		logger.Error("Error initializing folder %q. Error: %s", repositoryPath, err.Error())
		return false
	}

	return true
}
