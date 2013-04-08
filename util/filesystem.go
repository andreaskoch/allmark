// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"bufio"
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

func CreateDirectory(directoryPath string) bool {
	err := os.MkdirAll(directoryPath, 0700)
	return err == nil
}

func CreateFile(filePath string) (success bool, err error) {

	// make sure the parent directory exists
	directory := filepath.Dir(filePath)
	if !DirectoryExists(directory) {
		if !CreateDirectory(directory) {
			return false, fmt.Errorf("Cannot create the directory for the given file %q.", filePath)
		}
	}

	// create the file
	if _, err := os.Create(filePath); err != nil {
		return false, fmt.Errorf("Could not create file %q. Error: ", filePath, err)
	}

	return true, nil
}

func FileExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}

	file, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return !file.IsDir()
}

func DirectoryExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}

	file, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return file.IsDir()
}

// Gets the current working directory in which this application is being executed.
func GetWorkingDirectory() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "."
	}

	return workingDirectory
}
