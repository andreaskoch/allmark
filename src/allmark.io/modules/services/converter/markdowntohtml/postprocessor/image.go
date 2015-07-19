// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package postprocessor

import (
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/pattern"
	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter/markdowntohtml/common"
	"allmark.io/modules/services/converter/markdowntohtml/util"
	"fmt"
	"regexp"
	"strings"
)

var (
	// ![*image title (optional)*](*image path*)
	markdownPattern = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
)

func newImagePostprocessor(pathProvider paths.Pather, baseRoute route.Route, files []*model.File, imageProvider *common.ImageProvider) *imagePostProcessor {
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
	imageProvider *common.ImageProvider
}

func (postprocessor *imagePostProcessor) Convert(markdown string) (convertedContent string, converterError error) {

	convertedContent = markdown

	for {

		// search for files-extension code
		found, matches := pattern.IsMatch(convertedContent, markdownPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// normalize the path with the current path provider
		path = postprocessor.pathProvider.Path(path)

		// get the matching file
		file := postprocessor.getMatchingFile(path)
		if file == nil {

			// this is not an internal image reference
			imageCode := fmt.Sprintf(`<img src="%s" alt="%s"/>`, path, title)
			convertedContent = strings.Replace(convertedContent, originalText, imageCode, 1)
			continue
		}

		// get the image code
		imageCode := postprocessor.imageProvider.GetImageCodeWithLink(title, file.Route())

		// replace markdown with the image code
		convertedContent = strings.Replace(convertedContent, originalText, imageCode, 1)

	}

	return convertedContent, nil
}

func (postprocessor *imagePostProcessor) getMatchingFile(path string) *model.File {
	for _, file := range postprocessor.files {
		if file.Route().IsMatch(path) && util.IsImageFile(file) {
			return file
		}
	}

	return nil
}
