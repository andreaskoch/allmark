// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import (
	"github.com/russross/blackfriday"
	"path/filepath"
	"strings"
)

func ToHtml(markdown string) (html string) {
	return string(blackfriday.MarkdownCommon([]byte(markdown)))
}

func IsMarkdownFile(fileNameOrPath string) bool {
	fileExtension := strings.ToLower(filepath.Ext(fileNameOrPath))
	switch fileExtension {
	case ".md", ".markdown", ".mdown":
		return true
	default:
		return false
	}

	panic("Unreachable")
}
