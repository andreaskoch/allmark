// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/common/util/hashutil"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	ReservedDirectoryNames = []string{config.FilesDirectoryName, config.MetaDataFolderName}
)

func getHash(filepath string, route *route.Route) (string, error) {

	// fallback file hash
	fileHash, fallbackHashErr := getStringHash("")
	if fallbackHashErr != nil {
		return "", fallbackHashErr
	}

	// file hash
	if isFile, _ := fsutil.IsFile(filepath); isFile {
		if hash, err := getFileHash(filepath); err == nil {
			fileHash = hash
		}
	}

	// route hash
	routeHash, routeHashErr := getRouteHash(route)
	if routeHashErr != nil {
		return "", routeHashErr
	}

	// return the combined hash
	return fmt.Sprintf("%s+%s", routeHash, fileHash), nil
}

func getRouteHash(route *route.Route) (string, error) {
	return getStringHash(route.String())
}

func getStringHash(text string) (string, error) {
	routeReader := bytes.NewReader([]byte(text))
	return hashutil.GetHash(routeReader)
}

func getFileHash(path string) (string, error) {

	fileReader, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer fileReader.Close()

	return hashutil.GetHash(fileReader)
}

func getContent(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}

	defer file.Close()

	return ioutil.ReadAll(file)
}

func isReservedDirectory(path string) bool {

	if isFile, _ := fsutil.IsFile(path); isFile {
		path = filepath.Dir(path)
	}

	// get the directory name
	directoryName := strings.ToLower(filepath.Base(path))

	// all dot-directories are ignored
	if strings.HasPrefix(directoryName, ".") {
		return true
	}

	// check the reserved directory names
	for _, reservedDirectoryName := range ReservedDirectoryNames {
		if directoryName == strings.ToLower(reservedDirectoryName) {
			return true
		}
	}

	return false
}

func findMarkdownFileInDirectory(directory string) (found bool, file string) {
	entries, err := ioutil.ReadDir(directory)
	if err != nil {
		return false, ""
	}

	for _, element := range entries {

		if element.IsDir() {
			continue // skip directories
		}

		absoluteFilePath := filepath.Join(directory, element.Name())
		if isMarkdown := isMarkdownFile(absoluteFilePath); isMarkdown {
			return true, absoluteFilePath
		}
	}

	return false, ""
}

func getChildDirectories(directory string) []string {

	directories := make([]string, 0)
	directoryEntries, _ := ioutil.ReadDir(directory)
	for _, entry := range directoryEntries {

		if !entry.IsDir() {
			continue // skip files
		}

		childDirectory := filepath.Join(directory, entry.Name())
		if isReservedDirectory(childDirectory) {
			continue // skip reserved directories
		}

		// append directory
		directories = append(directories, childDirectory)
	}

	return directories
}

func isMarkdownFile(fileNameOrPath string) bool {
	fileExtension := strings.ToLower(filepath.Ext(fileNameOrPath))
	switch fileExtension {
	case ".md", ".markdown", ".mdown":
		return true
	default:
		return false
	}

	panic("Unreachable")
}
