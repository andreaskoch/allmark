// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/andreaskoch/allmark/config"
	"github.com/andreaskoch/allmark/renderer"
	"github.com/andreaskoch/allmark/server"
	"github.com/andreaskoch/allmark/util"
	"os"
	"strings"
)

const (
	CommandNameInit   = "init"
	CommandNameServe  = "serve"
	CommandNameRender = "render"
)

func main() {

	render := func(repositoryPath string) {
		config := config.GetConfig(repositoryPath)
		useTempDir := false
		renderer := renderer.New(repositoryPath, config, useTempDir)

		renderer.Execute()
	}

	serve := func(repositoryPath string) {

		config := config.GetConfig(repositoryPath)
		useTempDir := true
		server := server.New(repositoryPath, config, useTempDir)

		server.Serve()
	}

	init := func(repositoryPath string) {
		config.Initialize(repositoryPath)
	}

	parseCommandLineArguments(os.Args, func(commandName, repositoryPath string) (commandWasFound bool) {
		switch strings.ToLower(commandName) {
		case CommandNameInit:
			init(repositoryPath)

		case CommandNameRender:
			render(repositoryPath)

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

	} else {

		// use the current directory
		repositoryPath = util.GetWorkingDirectory()

	}

	// validate the supplied repository paths
	if !util.PathExists(repositoryPath) {
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
	fmt.Fprintf(os.Stderr, "  %7s  %s\n", CommandNameRender, "Start rendering the items in the specified repository")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Fork me on GitHub %q", "https://github.com/andreaskoch/allmark")

	os.Exit(2)
}
