// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"andyk/docs/indexer"
	"flag"
	"fmt"
	"os"
)

func main() {
	// define and parse application flags
	var repositoryPath = flag.String("repository", getWorkingDirectory(), "The path to a document repository.")

	flag.Usage = printUsageInformation
	flag.Parse()

	indexer.Index(*repositoryPath)
}

// Print usage information
var printUsageInformation = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

	flag.PrintDefaults()

	os.Exit(2)
}

// Gets the current working directory in which this application is being executed.
func getWorkingDirectory() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "."
	}

	return workingDirectory
}
