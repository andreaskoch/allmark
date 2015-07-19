// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	imageSrcPattern    = regexp.MustCompile(`src="([^"]+)"`)
	imageSrcSetPattern = regexp.MustCompile(`srcset="([^"]+)"`)
)

func lazyLoad(html string) string {

	html = lazyLoadSrcSet(html)
	html = lazyLoadSrc(html)

	return html
}

func lazyLoadSrc(html string) string {

	allMatches := imageSrcPattern.FindAllStringSubmatch(html, -1)
	for _, matches := range allMatches {

		if len(matches) != 2 {
			continue
		}

		// components
		originalText := strings.TrimSpace(matches[0])
		path := strings.TrimSpace(matches[1])

		// assemble the new link
		newLinkText := fmt.Sprintf(`data-sizes="auto" class="lazyload" data-src="%s"`, path)

		// replace the old text
		html = strings.Replace(html, originalText, newLinkText, -1)

	}

	return html
}

func lazyLoadSrcSet(html string) string {

	allMatches := imageSrcSetPattern.FindAllStringSubmatch(html, -1)
	for _, matches := range allMatches {

		if len(matches) != 2 {
			continue
		}

		// components
		originalText := strings.TrimSpace(matches[0])
		srcSetPaths := strings.TrimSpace(matches[1])

		// assemble the new link
		newLinkText := fmt.Sprintf(`data-sizes="auto" class="lazyload" data-srcset="%s"`, srcSetPaths)

		// replace the old text
		html = strings.Replace(html, originalText, newLinkText, -1)

	}

	return html
}
