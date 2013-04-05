// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/andreaskoch/allmark/renderer"
	"github.com/andreaskoch/allmark/server"
	"github.com/andreaskoch/allmark/util"
	"os"
	"strings"
	"time"
)

const (
	CommandNameServe  = "serve"
	CommandNameRender = "render"
)

func main() {

	// render callback
	render := func(repositoryPath string) {
		renderer.RenderRepository(repositoryPath)

		for {
			time.Sleep(100 * time.Millisecond)
		}
	}

	// serve callback
	serve := func(repositoryPath string) {
		server.Serve(repositoryPath)
	}

	parseCommandLineArguments(os.Args, render, serve)
}

func parseCommandLineArguments(args []string, renderCallback func(repositoryPath string), serveCallback func(repositoryPath string)) {

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
	_, err := util.IsValidDirectory(repositoryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "One or more of the supplied repository paths is invalid.\nError: %v", err)
		return
	}

	// Read the command parameter
	commandName := strings.ToLower(args[1])
	switch commandName {
	case CommandNameServe:
		{
			fmt.Printf("Serving repository: %v\n", repositoryPath)
			serveCallback(repositoryPath)
		}
	case CommandNameRender:
		{
			fmt.Printf("Rendering repository: %v\n", repositoryPath)
			renderCallback(repositoryPath)
		}

	default:
		{
			printUsageInformation(args)
		}
	}
}

// Print usage information
func printUsageInformation(args []string) {
	executeableName := args[0]

	fmt.Fprintf(os.Stderr, "%s - %s\n", executeableName, "A markdown web server and renderer")
	fmt.Fprintf(os.Stderr, "\nUsage:\n%s %s %s\n", executeableName, "<command>", "<repository path>")
	fmt.Fprintf(os.Stderr, "\nAvailable commands:\n")
	fmt.Fprintf(os.Stderr, "  %7s  %s\n", CommandNameServe, "Start serving the supplied repository via HTTP")
	fmt.Fprintf(os.Stderr, "  %7s  %s\n", CommandNameRender, "Start rendering the items in the specified repository")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Fork me on GitHub %q", "https://github.com/andreaskoch/allmark")

	os.Exit(2)
}
