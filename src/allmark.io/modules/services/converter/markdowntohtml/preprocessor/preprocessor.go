// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package preprocessor

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter/markdowntohtml/common"
	"allmark.io/modules/services/thumbnail"
)

// Preprocessor provides pre-processing capabilties for markdown code.
type Preprocessor struct {
	logger         logger.Logger
	thumbnailIndex *thumbnail.Index
}

// New creates an instance of a Markdown Preprocessor.
func New(logger logger.Logger, thumbnailIndex *thumbnail.Index) *Preprocessor {
	return &Preprocessor{
		logger:         logger,
		thumbnailIndex: thumbnailIndex,
	}
}

// Convert converts all markdown extensions in the supplied markdown to normal markdown code or HTML.
func (preprocessor *Preprocessor) Convert(
	aliasResolver func(alias string) *model.Item,
	rootPathProvider, itemContentPathProvider paths.Pather,
	itemRoute route.Route,
	files []*model.File,
	markdown string) (processedMarkdown string, errors error) {

	imageProvider := common.NewImageProvider(rootPathProvider, preprocessor.thumbnailIndex)

	// markdown extension: audio
	audioConverter := newAudioExtension(itemContentPathProvider, files)
	markdown, audioConversionError := audioConverter.Convert(markdown)
	if audioConversionError != nil {
		preprocessor.logger.Warn("Error while converting audio extensions. Error: %s", audioConversionError)
	}

	// markdown extension: video
	videoConverter := newVideoExtension(itemContentPathProvider, files)
	markdown, videoConversionError := videoConverter.Convert(markdown)
	if videoConversionError != nil {
		preprocessor.logger.Warn("Error while converting video extensions. Error: %s", videoConversionError)
	}

	// markdown extension: files
	filesConverter := newFilesExtension(itemContentPathProvider, itemRoute, files)
	markdown, filesConversionError := filesConverter.Convert(markdown)
	if filesConversionError != nil {
		preprocessor.logger.Warn("Error while converting files extensions. Error: %s", filesConversionError)
	}

	// markdown extension: filepreview
	filePreviewConverter := newFilePreviewExtension(itemContentPathProvider, files)
	markdown, filePreviewConversionError := filePreviewConverter.Convert(markdown)
	if filePreviewConversionError != nil {
		preprocessor.logger.Warn("Error while converting file preview extensions. Error: %s", filePreviewConversionError)
	}

	// markdown extension: imagegallery
	imagegalleryConverter := newImageGalleryExtension(itemContentPathProvider, itemRoute, files, imageProvider)
	markdown, imagegalleryConversionError := imagegalleryConverter.Convert(markdown)
	if imagegalleryConversionError != nil {
		preprocessor.logger.Warn("Error while converting image gallery extensions. Error: %s", imagegalleryConversionError)
	}

	// markdown extension: csv table
	csvTableConverter := newCSVExtension(itemContentPathProvider, files)
	markdown, csvTableConversionError := csvTableConverter.Convert(markdown)
	if csvTableConversionError != nil {
		preprocessor.logger.Warn("Error while converting csv table extensions. Error: %s", csvTableConversionError)
	}

	// markdown extension: reference
	referenceConverter := newReferenceExtension(itemContentPathProvider, aliasResolver)
	markdown, referenceConversionError := referenceConverter.Convert(markdown)
	if referenceConversionError != nil {
		preprocessor.logger.Warn("Error while converting reference extensions. Error: %s", referenceConversionError)
	}

	return markdown, nil

}
