// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package preprocessor

import (
	"github.com/elWyatt/allmark/common/logger"
	"github.com/elWyatt/allmark/common/paths"
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/model"
	"github.com/elWyatt/allmark/services/converter/markdowntohtml/imageprovider"
)

// Preprocessor provides pre-processing capabilties for markdown code.
type Preprocessor struct {
	logger        logger.Logger
	imageProvider *imageprovider.ImageProvider
}

// New creates an instance of a Markdown Preprocessor.
func New(logger logger.Logger, imageProvider *imageprovider.ImageProvider) *Preprocessor {
	return &Preprocessor{
		logger:        logger,
		imageProvider: imageProvider,
	}
}

// Convert converts all markdown extensions in the supplied markdown to normal markdown code or HTML.
func (preprocessor *Preprocessor) Convert(
	aliasResolver func(alias string) *model.Item,
	pathProvider paths.Pather,
	itemRoute route.Route,
	files []*model.File,
	markdown string) (processedMarkdown string, errors error) {

	// markdown extension: audio
	audioConverter := newAudioExtension(pathProvider, files)
	markdown, audioConversionError := audioConverter.Convert(markdown)
	if audioConversionError != nil {
		preprocessor.logger.Warn("Error while converting audio extensions. Error: %s", audioConversionError)
	}

	// markdown extension: video
	videoConverter := newVideoExtension(pathProvider, files)
	markdown, videoConversionError := videoConverter.Convert(markdown)
	if videoConversionError != nil {
		preprocessor.logger.Warn("Error while converting video extensions. Error: %s", videoConversionError)
	}

	// markdown extension: files
	filesConverter := newFilesExtension(pathProvider, itemRoute, files)
	markdown, filesConversionError := filesConverter.Convert(markdown)
	if filesConversionError != nil {
		preprocessor.logger.Warn("Error while converting files extensions. Error: %s", filesConversionError)
	}

	// markdown extension: filepreview
	filePreviewConverter := newFilePreviewExtension(pathProvider, files)
	markdown, filePreviewConversionError := filePreviewConverter.Convert(markdown)
	if filePreviewConversionError != nil {
		preprocessor.logger.Warn("Error while converting file preview extensions. Error: %s", filePreviewConversionError)
	}

	// markdown extension: imagegallery
	imagegalleryConverter := newImageGalleryExtension(pathProvider, itemRoute, files, preprocessor.imageProvider)
	markdown, imagegalleryConversionError := imagegalleryConverter.Convert(markdown)
	if imagegalleryConversionError != nil {
		preprocessor.logger.Warn("Error while converting image gallery extensions. Error: %s", imagegalleryConversionError)
	}

	// markdown extension: csv table
	csvTableConverter := newCSVExtension(pathProvider, files)
	markdown, csvTableConversionError := csvTableConverter.Convert(markdown)
	if csvTableConversionError != nil {
		preprocessor.logger.Warn("Error while converting csv table extensions. Error: %s", csvTableConversionError)
	}

	// markdown extension: reference
	referenceConverter := newReferenceExtension(pathProvider, aliasResolver)
	markdown, referenceConversionError := referenceConverter.Convert(markdown)
	if referenceConversionError != nil {
		preprocessor.logger.Warn("Error while converting reference extensions. Error: %s", referenceConversionError)
	}

	// markdown extension: mermaid
	mermaidConverter := newMermaidExtension(pathProvider, files)
	markdown, mermaidConversionError := mermaidConverter.Convert(markdown)
	if mermaidConversionError != nil {
		preprocessor.logger.Warn("Error while converting reference extensions. Error: %s", referenceConversionError)
	}

	return markdown, nil

}
