// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package postprocessor

import (
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter/markdowntohtml/imageprovider"
	"regexp"
	"strings"
)

var (
	imageSourcePattern = regexp.MustCompile(`src="([^"]+)"`)
)

func newImagePostprocessor(pathProvider paths.Pather, baseRoute route.Route, files []*model.File, imageProvider *imageprovider.ImageProvider) *imagePostProcessor {
	return &imagePostProcessor{
		pathProvider:  pathProvider,
		files:         files,
		imageProvider: imageProvider,
	}
}

type imagePostProcessor struct {
	pathProvider  paths.Pather
	base          route.Route
	files         []*model.File
	imageProvider *imageprovider.ImageProvider
}

func (postprocessor *imagePostProcessor) Convert(markdown string) (convertedContent string, converterError error) {

	convertedContent = markdown

	for _, match := range imageSourcePattern.FindAllStringSubmatch(convertedContent, -1) {

		if len(match) != 2 {
			continue
		}

		// parameters
		originalText := strings.TrimSpace(match[0])
		filePath := strings.TrimSpace(match[1])
		fileRoute := route.Combine(postprocessor.base, route.NewFromRequest(filePath))
		path := fileRoute.Value()

		// normalize the path with the current path provider
		path = postprocessor.pathProvider.Path(path)

		// get the matching file
		file := postprocessor.getMatchingFile(path)
		if file == nil {

			// this is not an internal image reference
			continue

		}

		// get the image path (src="...", srcset="...")
		imagePath := postprocessor.imageProvider.GetImagePath(postprocessor.pathProvider, file.Route())

		// replace markdown with the image code
		convertedContent = strings.Replace(convertedContent, originalText, imagePath, 1)

	}

	return convertedContent, nil
}

func (postprocessor *imagePostProcessor) getMatchingFile(path string) *model.File {
	for _, file := range postprocessor.files {
		if file.Route().IsMatch(path) && model.IsImageFile(file) {
			return file
		}
	}

	return nil
}
