// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/andreaskoch/allmark/common/config"
	"github.com/andreaskoch/allmark/common/logger/console"
	"github.com/andreaskoch/allmark/common/logger/loglevel"
	"github.com/andreaskoch/allmark/common/shutdown"
	"github.com/andreaskoch/allmark/common/util/fsutil"
	"github.com/andreaskoch/allmark/dataaccess/filesystem"
	"github.com/andreaskoch/allmark/services/initialization"
	"github.com/andreaskoch/allmark/services/parser"
	"github.com/andreaskoch/allmark/services/thumbnail"
	"github.com/andreaskoch/allmark/web/server"
	// "github.com/davecheney/profile"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

const (

	// CommandNameInit contains the name of the init action
	CommandNameInit = "init"

	// CommandNameServe contains the name of the serve action
	CommandNameServe = "serve"

	// CommandNameVersion contains the name of the version action
	CommandNameVersion = "version"
)

var version = "v0.10.0-dev"

var (
	serveFlags       = flag.NewFlagSet("serve-flags", flag.ContinueOnError)
	secure           = serveFlags.Bool("secure", false, "Use HTTPs only")
	logLevelOverride = serveFlags.String("loglevel", "", "Log level")
	reindex          = serveFlags.Bool("reindex", false, "Enable reindexing")
	livereload       = serveFlags.Bool("livereload", false, "Enable live-reload")
)

func main() {

	// defer profile.Start(profile.CPUProfile).Stop()

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

		case CommandNameVersion:
			printVersionInformation()
			return true

		default:
			return false
		}
	})
}

func parseCommandLineArguments(args []string, commandHandler func(commandName, repositoryPath string) (commandWasFound bool)) {

	remainingArguments := args

	// check if the mandatory amount of
	// command line parameters has been
	// supplied. If not, print usage information.
	if len(remainingArguments) < 2 {
		printUsageInformation(args)
		return
	}

	commandName := strings.ToLower(remainingArguments[1])
	remainingArguments = remainingArguments[2:]

	// Read the repository path parameters
	var repositoryPath string
	if len(remainingArguments) > 0 && !isCommandlineFlag(remainingArguments[0]) {

		// use supplied repository path
		repositoryPath = remainingArguments[0]
		remainingArguments = remainingArguments[1:]

		if isFile, _ := fsutil.IsFile(repositoryPath); isFile {
			repositoryPath = filepath.Dir(repositoryPath)
		}

	} else {

		// use the current directory
		repositoryPath = fsutil.GetWorkingDirectory()

	}

	// use the rest of the arguments to parse flags
	if len(remainingArguments) > 0 {
		serveFlags.Parse(remainingArguments)
	}

	// validate the supplied repository paths
	if !fsutil.PathExists(repositoryPath) {
		fmt.Fprintf(os.Stderr, "The specified repository paths %q is does not exist.", repositoryPath)
		return
	}

	// Read the command parameter and execute the command handler
	if commandWasFound := commandHandler(commandName, repositoryPath); !commandWasFound {
		printUsageInformation(args)
	}
}

// Print usage information
func printUsageInformation(args []string) {
	executeableName := args[0]

	fmt.Fprintf(os.Stderr, "%s - %s (Version: %s)\n", executeableName, "The standalone markdown webserver", version)
	fmt.Fprintf(os.Stderr, "\nUsage:\n%s %s %s\n", executeableName, "<command>", "<repository path>")
	fmt.Fprintf(os.Stderr, "\nAvailable commands:\n")
	fmt.Fprintf(os.Stderr, "  %7s  %s\n", CommandNameInit, "Initialize the configuration")
	fmt.Fprintf(os.Stderr, "  %7s  %s\n", CommandNameServe, "Start serving the supplied repository via HTTP and HTTPs")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Fork me on GitHub %q\n", "https://github.com/andreaskoch/allmark")

	os.Exit(2)
}

func serve(repositoryPath string) bool {

	// get the configuration
	configuration := config.Get(repositoryPath)

	// check if https shall be forced
	if *secure {
		configuration.Server.HTTPS.Enabled = true
		configuration.Server.HTTPS.Force = true
	}

	// check if indexing is enabled
	if *reindex {
		configuration.Indexing.Enabled = true
		configuration.Indexing.IntervalInSeconds = config.DefaultIndexingIntervalInSeconds
	}

	// check if live-reload is enabled
	if *livereload {
		configuration.LiveReload.Enabled = true
	}

	// create a logger
	logger := console.New(loglevel.FromString(configuration.LogLevel))
	if *logLevelOverride != "" {
		logger = console.New(loglevel.FromString(*logLevelOverride))
	}

	// data access
	repository, err := filesystem.NewRepository(logger, repositoryPath, *configuration)
	if err != nil {
		logger.Fatal("Unable to create a repository. Error: %s", err)
	}

	// thumbnail index
	thumbnailIndex := thumbnail.EmptyIndex()
	if configuration.Conversion.Thumbnails.Enabled {

		thumbnailIndexFilePath := configuration.ThumbnailIndexFilePath()
		thumbnailFolder := configuration.ThumbnailFolder()

		if !fsutil.CreateDirectory(thumbnailFolder) {
			logger.Fatal("Could not create the thumbnail folder %q", thumbnailFolder)
		}

		thumbnailIndex = thumbnail.NewIndex(logger, thumbnailIndexFilePath, thumbnailFolder)

		// thumbnail conversion service
		thumbnail.NewConversionService(logger, repository, thumbnailIndex)

	}

	// parser
	itemParser, err := parser.New(logger)
	if err != nil {
		logger.Fatal("Unable to instantiate a parser. Error: %s", err)
	}

	// server
	server, err := server.New(logger, *configuration, repository, itemParser, thumbnailIndex)
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

func printVersionInformation() {
	fmt.Println(version)
}

func isCommandlineFlag(argument string) bool {
	return strings.HasPrefix(argument, "-")
}
