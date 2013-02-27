// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"andyk/docs/renderer"
	"andyk/docs/server"
	"andyk/docs/util"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	CommandNameServe  = "serve"
	CommandNameRender = "render"
)

func main() {

	// render callback
	render := func(repositoryPaths []string) {
		renderer.Render(repositoryPaths)
	}

	// serve callback
	serve := func(repositoryPaths []string) {
		server.Serve(repositoryPaths)
	}

	parseCommandLineArguments(os.Args, render, serve)
}

func parseCommandLineArguments(args []string, renderCallback func(repositoryPaths []string), serveCallback func(repositoryPaths []string)) {

	// check if the mandatory amount of
	// command line parameters has been
	// supplied. If not, print usage information.
	if len(args) < 2 {
		printUsageInformation(args)
		return
	}

	// Read the repository path parameters
	repositoryPaths := make([]string, 1, 1)
	if len(args) > 2 {

		// use supplied repository paths
		repositoryPaths = args[2:]

	} else {

		// use the current directory
		repositoryPaths[0] = getWorkingDirectory()

	}

	// validate the supplied repository paths
	_, err := repositoryPathsAreValid(repositoryPaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "One or more of the supplied repository paths is invalid.\nError: %v", err)
		return
	}

	// Read the command parameter
	commandName := strings.ToLower(args[1])
	switch commandName {
	case CommandNameServe:
		{
			fmt.Printf("Serving repositories: %v\n", strings.Join(repositoryPaths, ", "))
			serveCallback(repositoryPaths)
		}
	case CommandNameRender:
		{
			fmt.Printf("Rendering repositories: %v\n", strings.Join(repositoryPaths, ", "))
			renderCallback(repositoryPaths)
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

	fmt.Fprintf(os.Stderr, "%s - %s\n", executeableName, "A markdown server and and renderer.")
	fmt.Fprintf(os.Stderr, "\nUsage:\n%s %s %s %s %s\n", executeableName, "<command>", "[<repository path>]", "[<...>]", "[<repository path>]")
	fmt.Fprintf(os.Stderr, "\nAvailable commands:\n")
	fmt.Fprintf(os.Stderr, "  %s    %s\n", CommandNameServe, "Serve the supplied repositor(y/ies) via HTTP.")
	fmt.Fprintf(os.Stderr, "  %s   %s\n", CommandNameRender, "Render the items in the supplied repositor(y/ies).")

	os.Exit(2)
}

func repositoryPathsAreValid(repositoryPaths []string) (bool, error) {

	pathCounter := make(map[string]int)

	for _, path := range repositoryPaths {

		// A repository path cannot be empty
		if strings.TrimSpace(path) == "" {
			return false, errors.New("A repository path cannot be empty.")
		}

		// Get the absolute file path
		absoluteFilePath, absoluteFilePathError := filepath.Abs(path)
		if absoluteFilePathError != nil {
			return false, errors.New(fmt.Sprintf("Cannot determine the absolute repository path for the supplied repository: %v", path))
		}

		// The respository path must be accessible
		if !util.FileExists(absoluteFilePath) {
			return false, errors.New(fmt.Sprintf("The repository path \"%s\" cannot be accessed.", path))
		}

		// Check for duplicates
		normalizedPath := strings.ToLower(absoluteFilePath)
		pathCounter[normalizedPath] += 1
		if pathCounter[normalizedPath] > 1 {
			return false, errors.New(fmt.Sprintf("The repository paths cannot contain the same path twice: %s", path))
		}
	}

	return true, nil
}

// Gets the current working directory in which this application is being executed.
func getWorkingDirectory() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "."
	}

	return workingDirectory
}
