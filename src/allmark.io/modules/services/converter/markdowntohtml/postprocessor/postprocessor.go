// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package postprocessor

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter/markdowntohtml/common"
	"allmark.io/modules/services/thumbnail"
)

// Postprocessor provides post-processing capabilities for HTML code.
type Postprocessor struct {
	logger         logger.Logger
	thumbnailIndex *thumbnail.Index
}

// New creates a new Postprocessor.
func New(logger logger.Logger, thumbnailIndex *thumbnail.Index) *Postprocessor {
	return &Postprocessor{
		logger:         logger,
		thumbnailIndex: thumbnailIndex,
	}
}

// Convert applies post-processing to the supplied HTML code.
func (postprocessor *Postprocessor) Convert(
	rootPathProvider, itemContentPathProvider paths.Pather,
	itemRoute route.Route,
	files []*model.File,
	html string) (convertedContent string, converterError error) {

	// Thumbnails
	imageProvider := common.NewImageProvider(rootPathProvider, postprocessor.thumbnailIndex)
	imagePostProcessor := newImagePostprocessor(itemContentPathProvider, itemRoute, files, imageProvider)
	html, imageConversionError := imagePostProcessor.Convert(html)
	if imageConversionError != nil {
		postprocessor.logger.Warn("Error while converting images/thumbnails. Error: %s", imageConversionError)
	}

	// Rewrite Links
	html = rewireLinks(itemContentPathProvider, files, html)

	// Add Emojis
	html = addEmojis(html)

	return html, nil
}
