// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"github.com/andreaskoch/allmark/markdown"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
)

func renderMarkdown(fileIndex *repository.FileIndex, pathProvider *path.Provider, markdownCode string) string {
	return markdown.ToHtml(markdownCode)
}
