// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filesystem

import (
	"bytes"
	"fmt"
	"github.com/andreaskoch/allmark2/common/content"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/util/fsutil"
	"github.com/andreaskoch/allmark2/common/util/hashutil"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func newContentProvider(path string, route *route.Route) *content.ContentProvider {

	// mimeType
	mimeType := func() (string, error) {
		return getMimeType(path)
	}

	// content provider
	dataProvider := func() ([]byte, error) {
		return getData(path)
	}

	// hash provider
	hashProvider := func() (string, error) {

		// file hash
		fileHash, fileHashErr := getHash(path, route)
		if fileHashErr != nil {
			return "", fmt.Errorf("Unable to determine the hash for file %q. Error: %s", path, fileHashErr)
		}

		return fileHash, nil
	}

	// last modified provider
	lastModifiedProvider := func() (time.Time, error) {
		return fsutil.GetModificationTime(path)
	}

	return content.NewProvider(mimeType, dataProvider, hashProvider, lastModifiedProvider)
}

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
