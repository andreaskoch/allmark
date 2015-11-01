// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package postprocessor

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter/markdowntohtml/imageprovider"
)

// Postprocessor provides post-processing capabilities for HTML code.
type Postprocessor struct {
	logger        logger.Logger
	imageProvider *imageprovider.ImageProvider
}

// New creates a new Postprocessor.
func New(logger logger.Logger, imageProvider *imageprovider.ImageProvider) *Postprocessor {
	return &Postprocessor{
		logger:        logger,
		imageProvider: imageProvider,
	}
}

// Convert applies post-processing to the supplied HTML code.
func (postprocessor *Postprocessor) Convert(
	pathProvider paths.Pather,
	itemRoute route.Route,
	files []*model.File,
	html string) (convertedContent string, converterError error) {

	// Thumbnails
	imagePostProcessor := newImagePostprocessor(pathProvider, itemRoute, files, postprocessor.imageProvider)
	html, imageConversionError := imagePostProcessor.Convert(html)
	if imageConversionError != nil {
		postprocessor.logger.Warn("Error while converting images/thumbnails. Error: %s", imageConversionError)
	}

	// Rewrite Links
	html = rewireLinks(pathProvider, itemRoute, files, html)

	// Add Emojis
	html = addEmojis(html)

	return html, nil
}
