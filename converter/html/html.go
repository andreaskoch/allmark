// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"github.com/andreaskoch/allmark/repository"
	"strings"
)

func ToHtml(item *repository.Item) *repository.Item {

	// assign the raw markdown content for the add-ins to work on
	item.ConvertedContent = strings.TrimSpace(strings.Join(item.RawContent, "\n"))

	// render markdown extensions
	renderImageGalleries(item)
	renderFileLinks(item)
	renderCSVTables(item)
	renderPDFs(item)
	renderVideos(item)
	renderAudio(item)

	// render markdown
	renderMarkdown(item)

	return item
}
