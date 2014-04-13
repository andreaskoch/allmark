// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdowntohtml

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/audio"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/csvtable"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/filepreview"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/files"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/imagegallery"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/markdown"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/pdf"
	"github.com/andreaskoch/allmark2/services/conversion/markdowntohtml/video"
	"regexp"
	"strings"
)

var (
	// [*description text*](*folder path*)
	markdownLinkPattern = regexp.MustCompile(`\[(.*)\]\(([^)]+)\)`)

	markdownItemLinkPattern = regexp.MustCompile(`\[(.*)\]\(/([^)]+)\)`)
)

type Converter struct {
	logger logger.Logger
}

func New(logger logger.Logger) (*Converter, error) {
	return &Converter{
		logger: logger,
	}, nil
}

// Convert the supplied item with all paths relative to the supplied base route
func (converter *Converter) Convert(pathProvider paths.Pather, item *model.Item) (convertedContent string, conversionError error) {

	converter.logger.Debug("Converting item %q.", item)

	content := item.Content

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
	filesConverter := files.New(pathProvider, item.Files())
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
	imagegalleryConverter := imagegallery.New(pathProvider, item.Files())
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

	// fix links
	content = rewireLinks(pathProvider, item, content)

	// markdown to html
	content = markdown.Convert(content)

	return content, nil
}

func rewireLinks(pathProvider paths.Pather, item *model.Item, markdown string) string {

	allMatches := markdownLinkPattern.FindAllStringSubmatch(markdown, -1)
	for _, matches := range allMatches {

		if len(matches) != 3 {
			continue
		}

		// components
		originalText := strings.TrimSpace(matches[0])
		descriptionText := strings.TrimSpace(matches[1])
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
			continue
		}

		matchingFilePath := pathProvider.Path(fullFileRoute.Value())

		// assemble the new link
		newLinkText := fmt.Sprintf("[%s](%s)", descriptionText, matchingFilePath)

		// replace the old text
		markdown = strings.Replace(markdown, originalText, newLinkText, 1)

	}

	return markdown
}

func getMatchingFiles(path string, item *model.Item) *model.File {
	for _, file := range item.Files() {
		if file.Route().IsMatch(path) {
			return file
		}
	}

	return nil
}
