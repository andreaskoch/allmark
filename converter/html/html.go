// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/types"
)

func Convert(item *repository.Item, pathProvider *path.Provider) string {

	// files
	files := item.Files

	// assign the raw markdown content for the add-ins to work on
	convertedContent := item.RawContent

	// render markdown extensions
	convertedContent = renderImageGalleries(files, pathProvider, convertedContent)
	convertedContent = renderFileLinks(files, pathProvider, convertedContent)
	convertedContent = renderCSVTables(files, pathProvider, convertedContent)
	convertedContent = renderPDFs(files, pathProvider, convertedContent)
	convertedContent = renderVideos(files, pathProvider, convertedContent)
	convertedContent = renderAudio(files, pathProvider, convertedContent)

	// render markdown
	convertedContent = renderMarkdown(files, pathProvider, convertedContent)

	switch itemType := item.MetaData.ItemType; itemType {
	case types.PresentationItemType:
		convertedContent = renderPresentation(convertedContent)
	}

	return convertedContent
}
