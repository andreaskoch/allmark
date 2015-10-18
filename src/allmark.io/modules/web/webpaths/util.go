// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"regexp"
)

var (
	// A pattern matching prococol prefixes (e.g. http://, https://, ftp://, bitcoin:, mailto: and any other)
	protocolPrefixPattern = regexp.MustCompile(`^\w+:`)
)

// IsAbsoluteURI checks if the given uri is absolute or not.
func IsAbsoluteURI(uri string) bool {
	uriHasProtocolPrefix := protocolPrefixPattern.MatchString(uri)
	return uriHasProtocolPrefix
}
