// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"github.com/russross/blackfriday"
)

func markdownToHtml(markdown string) (html string) {
	return string(blackfriday.MarkdownCommon([]byte(markdown)))
}
