// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"strings"
)

func CombineUrlComponents(baseUrl string, components ...string) string {
	url := StripTrailingUrlDirectorySeperator(baseUrl)

	for _, component := range components {
		url += UrlDirectorySeperator + StripTrailingUrlDirectorySeperator(component)
	}

	return url
}

func StripTrailingUrlDirectorySeperator(urlComponent string) string {

	url := urlComponent
	for strings.LastIndex(url, UrlDirectorySeperator)+1 == len(url) && len(url) != 0 {
		url = strings.TrimRight(url, UrlDirectorySeperator)
	}

	return url
}

func StripLeadingUrlDirectorySeperator(urlComponent string) string {

	url := urlComponent

	for strings.Index(url, UrlDirectorySeperator) == 0 {
		url = strings.TrimLeft(url, UrlDirectorySeperator)
	}

	return url
}
