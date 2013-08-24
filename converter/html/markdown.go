// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"github.com/andreaskoch/allmark/markdown"
	"github.com/andreaskoch/allmark/repository"
)

func renderMarkdown(item *repository.Item, rawContent string) string {
	return markdown.ToHtml(rawContent)
}
