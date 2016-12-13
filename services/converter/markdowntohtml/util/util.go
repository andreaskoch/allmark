// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"strings"
)

func GetHtmlLinkCode(title, path string) string {
	return fmt.Sprintf(`<a href="%s" target="_blank" title="%s">%s</a>`, path, title, title)
}

// IsInternalLink returns true if the supplied link is an internal or an external link.
func IsInternalLink(link string) bool {
	return !IsExternalLink(link)
}

func IsExternalLink(link string) bool {
	lowercase := strings.TrimSpace(strings.ToLower(link))
	isHttpLink := strings.HasPrefix(lowercase, "http:")
	isHttpsLink := strings.HasPrefix(lowercase, "https:")
	return isHttpLink || isHttpsLink
}
