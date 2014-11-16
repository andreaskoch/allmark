// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/pattern"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/converter/markdowntohtml/common"
	"github.com/andreaskoch/allmark2/services/converter/markdowntohtml/util"
	"regexp"
	"strings"
)

var (
	// ![*image title (optional)*](*image path*)
	markdownPattern = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
)

func New(pathProvider paths.Pather, baseRoute route.Route, files []*model.File, imageProvider *common.ImageProvider) *ImageExtension {
	return &ImageExtension{
		pathProvider:  pathProvider,
		files:         files,
		imageProvider: imageProvider,
	}
}

type ImageExtension struct {
	pathProvider  paths.Pather
	base          route.Route
	files         []*model.File
	imageProvider *common.ImageProvider
}

func (converter *ImageExtension) Convert(markdown string) (convertedContent string, converterError error) {

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
		path = converter.pathProvider.Path(path)

		// get the matching file
		file := converter.getMatchingFile(path)
		if file == nil {

			// this is not an internal image reference
			imageCode := fmt.Sprintf(`<img src="%s" alt="%s"/>`, path, title)
			convertedContent = strings.Replace(convertedContent, originalText, imageCode, 1)
			continue
		}

		// get the image code
		imageCode := converter.imageProvider.GetImageCodeWithLink(title, file.Route())

		// replace markdown with the image code
		convertedContent = strings.Replace(convertedContent, originalText, imageCode, 1)

	}

	return convertedContent, nil
}

func (converter *ImageExtension) getMatchingFile(path string) *model.File {
	for _, file := range converter.files {
		if file.Route().IsMatch(path) && util.IsImageFile(file) {
			return file
		}
	}

	return nil
}
