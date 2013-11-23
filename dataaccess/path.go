// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"regexp"
	"strings"
)

var (
	regexpWhitespacePattern         = regexp.MustCompile(`\s+`)
	regexpBackSlashPattern          = regexp.MustCompile(`\\+`)
	regexpdoubleForwardSlashPattern = regexp.MustCompile(`/+`)
)

// Normalize the supplied path to be used for an Item or File
func NormalizePath(path string) string {

	// trim spaces
	path = strings.TrimSpace(path)

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

	return path
}
