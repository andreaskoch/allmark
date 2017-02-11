// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fsutil

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func CreateDirectory(path string) bool {
	err := os.MkdirAll(path, 0700)
	return err == nil
}

func OpenFile(filepath string) (*os.File, error) {
	if !FileExists(filepath) {
		return nil, fmt.Errorf("The file %q does not exist.", filepath)
	}

	return os.OpenFile(filepath, os.O_RDWR, 0644)
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

func PathExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}

	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}

	return true
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

func IsFile(path string) (bool, error) {

	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir() == false, nil
}

func IsDirectory(path string) (bool, error) {

	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

// Gets the current working directory in which this application is being executed.
func GetWorkingDirectory() string {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "."
	}

	return workingDirectory
}

func GetSubDirectories(path string) []string {

	directories := make([]string, 0)

	if ok, _ := IsDirectory(path); !ok {
		fmt.Errorf("%q is not a directory.\n", path)
		return directories
	}

	directoryEntries, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Errorf("Cannot read directory %q.\n", path)
		return directories
	}

	for _, entry := range directoryEntries {
		if !entry.IsDir() {
			continue // skip files
		}

		subDirectory := filepath.Join(path, entry.Name())
		directories = append(directories, subDirectory)
	}

	return directories
}

func GetModificationTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		var t time.Time
		return t, err
	}

	return info.ModTime(), nil
}

// GetTempDirectory returns the path to a new temparory directory.
func GetTempDirectory() string {
	randomString := rand_str()
	tempDir := filepath.Join(os.TempDir(), randomString)

	if !CreateDirectory(tempDir) {
		panic(fmt.Sprintf("Cannot create temp directory %q.", tempDir))
	}

	return tempDir
}

func rand_str() string {
	size := 8
	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		panic(err)
	}

	rs := base64.URLEncoding.EncodeToString(rb)
	return rs
}
