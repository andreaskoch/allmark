// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdowntohtml

import (
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter/filetreerenderer"
	"allmark.io/modules/services/converter/markdowntohtml/audio"
	"allmark.io/modules/services/converter/markdowntohtml/common"
	"allmark.io/modules/services/converter/markdowntohtml/csvtable"
	"allmark.io/modules/services/converter/markdowntohtml/filepreview"
	"allmark.io/modules/services/converter/markdowntohtml/files"
	"allmark.io/modules/services/converter/markdowntohtml/image"
	"allmark.io/modules/services/converter/markdowntohtml/imagegallery"
	"allmark.io/modules/services/converter/markdowntohtml/markdown"
	"allmark.io/modules/services/converter/markdowntohtml/pdf"
	"allmark.io/modules/services/converter/markdowntohtml/reference"
	"allmark.io/modules/services/converter/markdowntohtml/video"
	"allmark.io/modules/services/thumbnail"
	"fmt"
	"regexp"
	"strings"
)

var (
	imageSrcPattern    = regexp.MustCompile(`src="([^"]+)"`)
	imageSrcSetPattern = regexp.MustCompile(`srcset="([^"]+)"`)

	htmlLinkPattern = regexp.MustCompile(`(src|href)="([^"]+)"`)
)

type Converter struct {
	logger         logger.Logger
	thumbnailIndex *thumbnail.Index
}

func New(logger logger.Logger, thumbnailIndex *thumbnail.Index) *Converter {
	return &Converter{
		logger:         logger,
		thumbnailIndex: thumbnailIndex,
	}
}

// Convert the supplied item with all paths relative to the supplied base route
func (converter *Converter) Convert(aliasResolver func(alias string) *model.Item, pathProvider paths.Pather, item *model.Item) (convertedContent string, converterError error) {

	converter.logger.Debug("Converting item %q.", item)

	itemRoute := item.Route()
	content := item.Content
	imageProvider := common.NewImageProvider(pathProvider, converter.thumbnailIndex)

	// markdown extension: image/thumbnails
	imageConverter := image.New(pathProvider, itemRoute, item.Files(), imageProvider)
	content, imageConversionError := imageConverter.Convert(content)
	if imageConversionError != nil {
		converter.logger.Warn("Error while converting images/thumbnails. Error: %s", imageConversionError)
	}

	// markdown extension: audio
	audioConverter := audio.New(pathProvider, item.Files())
	content, audioConversionError := audioConverter.Convert(content)
	if audioConversionError != nil {
		converter.logger.Warn("Error while converting audio extensions. Error: %s", audioConversionError)
	}

	// markdown extension: video
	videoConverter := video.New(pathProvider, item.Files())
	content, videoConversionError := videoConverter.Convert(content)
	if videoConversionError != nil {
		converter.logger.Warn("Error while converting video extensions. Error: %s", videoConversionError)
	}

	// markdown extension: pdf
	pdfConverter := pdf.New(pathProvider, item.Files())
	content, pdfConversionError := pdfConverter.Convert(content)
	if pdfConversionError != nil {
		converter.logger.Warn("Error while converting pdf extensions. Error: %s", pdfConversionError)
	}

	// markdown extension: files
	filesConverter := files.New(pathProvider, itemRoute, item.Files())
	content, filesConversionError := filesConverter.Convert(content)
	if filesConversionError != nil {
		converter.logger.Warn("Error while converting files extensions. Error: %s", filesConversionError)
	}

	// markdown extension: filepreview
	filePreviewConverter := filepreview.New(pathProvider, item.Files())
	content, filePreviewConversionError := filePreviewConverter.Convert(content)
	if filePreviewConversionError != nil {
		converter.logger.Warn("Error while converting file preview extensions. Error: %s", filePreviewConversionError)
	}

	// markdown extension: imagegallery
	imagegalleryConverter := imagegallery.New(pathProvider, itemRoute, item.Files(), imageProvider)
	content, imagegalleryConversionError := imagegalleryConverter.Convert(content)
	if imagegalleryConversionError != nil {
		converter.logger.Warn("Error while converting image gallery extensions. Error: %s", imagegalleryConversionError)
	}

	// markdown extension: csv table
	csvTableConverter := csvtable.New(pathProvider, item.Files())
	content, csvTableConversionError := csvTableConverter.Convert(content)
	if csvTableConversionError != nil {
		converter.logger.Warn("Error while converting csv table extensions. Error: %s", csvTableConversionError)
	}

	// markdown extension: reference
	referenceConverter := reference.New(pathProvider, aliasResolver)
	content, referenceConversionError := referenceConverter.Convert(content)
	if referenceConversionError != nil {
		converter.logger.Warn("Error while converting reference extensions. Error: %s", referenceConversionError)
	}

	// markdown to html
	content = markdown.Convert(content)

	// fix links
	content = converter.rewireLinks(pathProvider, item, content)

	// append the file list
	if item.IsFileCollection() && len(item.Files()) > 0 {
		fileTreeRenderer := filetreerenderer.New(pathProvider, itemRoute, item.Files())
		fileBaseFolder := getBaseFolder(itemRoute, item.Files())
		content += fileTreeRenderer.Render("Attachments", "attachments", fileBaseFolder)
	}

	return content, nil
}

func getBaseFolder(referenceRoute route.Route, files []*model.File) string {

	baseFolder := ""

	for _, file := range files {
		partialRoute := route.Intersect(referenceRoute, file.Route())

		if baseFolder == "" {
			baseFolder = partialRoute.FirstComponentName()
			continue
		}

		// abort if the base folders differ
		if baseFolder != partialRoute.FirstComponentName() {
			return ""
		}
	}

	return baseFolder
}

func (converter *Converter) rewireLinks(pathProvider paths.Pather, item *model.Item, html string) string {

	allMatches := htmlLinkPattern.FindAllStringSubmatch(html, -1)
	for _, matches := range allMatches {

		if len(matches) != 3 {
			continue
		}

		// components
		originalText := strings.TrimSpace(matches[0])
		linkType := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get matching file
		matchingFile := getMatchingFiles(path, item)

		// skip if no matching files are found
		if matchingFile == nil {
			continue
		}

		// assemble the new link path
		fullFileRoute, err := route.Combine(matchingFile.Parent(), matchingFile.Route())
		if err != nil {
			converter.logger.Error("%s", err)
			continue
		}

		matchingFilePath := pathProvider.Path(fullFileRoute.Value())

		// assemble the new link
		newLinkText := fmt.Sprintf("%s=\"%s\"", linkType, matchingFilePath)

		// replace the old text
		html = strings.Replace(html, originalText, newLinkText, -1)

	}

	return html
}

func getMatchingFiles(path string, item *model.Item) *model.File {
	for _, file := range item.Files() {
		if file.Route().IsMatch(path) {
			return file
		}
	}

	return nil
}

func LazyLoad(html string) string {

	html = lazyLoadSrcSet(html)
	html = lazyLoadSrc(html)

	return html
}

func lazyLoadSrc(html string) string {

	allMatches := imageSrcPattern.FindAllStringSubmatch(html, -1)
	for _, matches := range allMatches {

		if len(matches) != 2 {
			continue
		}

		// components
		originalText := strings.TrimSpace(matches[0])
		path := strings.TrimSpace(matches[1])

		// assemble the new link
		newLinkText := fmt.Sprintf("data-lazyload=\"true\" %s=\"%s\"", "data-src", path)

		// replace the old text
		html = strings.Replace(html, originalText, newLinkText, -1)

	}

	return html
}

func lazyLoadSrcSet(html string) string {

	allMatches := imageSrcSetPattern.FindAllStringSubmatch(html, -1)
	for _, matches := range allMatches {

		if len(matches) != 2 {
			continue
		}

		// components
		originalText := strings.TrimSpace(matches[0])
		srcSetPaths := strings.TrimSpace(matches[1])

		// assemble the new link
		newLinkText := fmt.Sprintf("data-lazyload=\"true\" %s=\"%s\"", "data-srcset", srcSetPaths)

		// replace the old text
		html = strings.Replace(html, originalText, newLinkText, -1)

	}

	return html
}
