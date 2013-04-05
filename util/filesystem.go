// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// readLine returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func readLine(bufferedReader *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = bufferedReader.ReadLine()
		ln = append(ln, line...)
	}

	return string(ln), err
}

// Get all lines of a given file
func GetLines(inFile io.Reader) []string {

	lines := make([]string, 0, 10)
	bufferedReader := bufio.NewReader(inFile)
	line, err := readLine(bufferedReader)
	for err == nil {
		lines = append(lines, line)
		line, err = readLine(bufferedReader)
	}

	return lines
}

func FileExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}

	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func IsValidDirectory(path string) (bool, error) {

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
	if !FileExists(absoluteFilePath) {
		return false, errors.New(fmt.Sprintf("The repository path \"%s\" cannot be accessed.", path))
	}

	return true, nil
}

// Gets the current working directory in which this application is being executed.
func GetWorkingDirectory() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "."
	}

	return workingDirectory
}
