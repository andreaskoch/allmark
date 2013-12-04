// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	regexpWhitespacePattern         = regexp.MustCompile(`\s+`)
	regexpBackSlashPattern          = regexp.MustCompile(`\\+`)
	regexpdoubleForwardSlashPattern = regexp.MustCompile(`/+`)
)

type Route struct {
	value string
}

func New(repositoryPath, itemPath string) (*Route, error) {

	// normalize the repository path
	normalizedRepositoryPath, err := normalizePath(repositoryPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to normalize the supplied repository path %q. Error: %s", repositoryPath, err)
	}

	// normalize the item path
	normalizedItemPath, err := normalizePath(itemPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to normalize the supplied item path %q. Error: %s", itemPath, err)
	}

	// strip the repository path from the item path
	routeValue := strings.Replace(normalizedItemPath, normalizedRepositoryPath, "", 1)
	return &Route{routeValue}, nil
}

func (route *Route) String() string {
	return route.value
}

// Normalize the supplied path to be used for an Item or File
func normalizePath(path string) (string, error) {

	// trim spaces
	path = strings.TrimSpace(path)

	// check if the path is empty
	if path == "" {
		return path, fmt.Errorf("A path cannot be empty.")
	}

	// replace all backslashes with a (single) forward slash
	path = regexpBackSlashPattern.ReplaceAllString(path, "/")

	// replace multiple forward slashes with a single forward slash
	path = regexpdoubleForwardSlashPattern.ReplaceAllString(path, "/")

	// remove leading slashes
	path = strings.TrimLeft(path, "/")

	// remove trailing slashes
	path = strings.TrimRight(path, "/")

	// replace duplicate spaces with a (single) url safe character
	path = regexpWhitespacePattern.ReplaceAllString(path, "+")

	return path, nil
}
