// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/content"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/common/util/hashutil"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func newFileContentProvider(path string, route route.Route) *content.ContentProvider {

	// mimeType
	mimeType := func() (string, error) {
		return getMimeType(path)
	}

	// content provider
	dataProvider := func(callback func(content io.ReadSeeker) error) error {

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		return callback(file)
	}

	// hash provider
	hashProvider := func() (string, error) {

		fileHash, fileHashErr := getHashFromFile(path, route)
		if fileHashErr != nil {
			return "", fmt.Errorf("Unable to determine the hash for file %q. Error: %s", path, fileHashErr)
		}

		return fileHash, nil
	}

	// last modified provider
	lastModifiedProvider := func() (time.Time, error) {
		return fsutil.GetModificationTime(path)
	}

	return content.NewContentProvider(mimeType,
		dataProvider,
		hashProvider,
		lastModifiedProvider)
}

func newFileContentProviderWithoutChecksum(path string, route route.Route) *content.ContentProvider {

	// mimeType
	mimeType := func() (string, error) {
		return getMimeType(path)
	}

	// content provider
	dataProvider := func(callback func(content io.ReadSeeker) error) error {

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer file.Close()

		return callback(file)
	}

	// hash provider
	hashProvider := func() (string, error) {

		// get the mod time
		modTime, modTimeErr := fsutil.GetModificationTime(path)
		if modTimeErr != nil {
			return "", modTimeErr
		}

		hashSource := fmt.Sprintf("%s - %s", route.Value(), modTime)
		return getStringHash(hashSource)
	}

	// last modified provider
	lastModifiedProvider := func() (time.Time, error) {
		return fsutil.GetModificationTime(path)
	}

	return content.NewContentProvider(mimeType,
		dataProvider,
		hashProvider,
		lastModifiedProvider)
}

func newTextContentProvider(text string, route route.Route) *content.ContentProvider {

	// mimeType
	mimeType := func() (string, error) {
		return "text/x-markdown", nil
	}

	// content provider
	dataProvider := func(callback func(content io.ReadSeeker) error) error {
		contentReader := strings.NewReader(text)
		return callback(contentReader)
	}

	// hash provider
	hashProvider := func() (string, error) {

		routeHash, routeHashErr := getStringHash(route.String())
		if routeHashErr != nil {
			return "", fmt.Errorf("Unable to determine the hash for route %q. Error: %s", route, routeHashErr)
		}

		contentHash, contentHashErr := getStringHash(text)
		if contentHashErr != nil {
			return "", fmt.Errorf("Unable to determine the hash for content %q. Error: %s", text, contentHashErr)
		}

		return fmt.Sprintf("%s%s", routeHash, contentHash), nil
	}

	// last modified provider
	lastModifiedProvider := func() (time.Time, error) {
		return time.Time{}, nil
	}

	return content.NewContentProvider(mimeType, dataProvider, hashProvider, lastModifiedProvider)
}

func getHashFromFile(filepath string, route route.Route) (string, error) {

	// fallback file hash
	fileHash, fallbackHashErr := getStringHash(filepath)
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
	return fmt.Sprintf("%s%s", routeHash, fileHash), nil
}

func getRouteHash(route route.Route) (string, error) {
	return getStringHash(route.String())
}

func getStringHash(text string) (string, error) {
	routeReader := strings.NewReader(text)
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

func getMimeType(path string) (string, error) {

	// content type detection
	// derive content type from file extension
	fileExtension := filepath.Ext(path)
	contentType := mime.TypeByExtension(fileExtension)
	if contentType == "" {
		// fallback: derive content type from data
		data, err := getData(path)
		if err != nil {
			return "", err
		}

		contentType = http.DetectContentType(data)
	}

	return contentType, nil
}

func getData(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}

	defer file.Close()

	return ioutil.ReadAll(file)
}
