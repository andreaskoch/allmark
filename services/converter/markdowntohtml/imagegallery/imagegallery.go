// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imagegallery

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/converter/markdowntohtml/pattern"
	"github.com/andreaskoch/allmark2/services/converter/markdowntohtml/util"
	"regexp"
	"strings"
)

var (
	// imagegallery: [*description text*](*folder path*)
	markdownPattern = regexp.MustCompile(`imagegallery: \[([^\]]+)\]\(([^)]+)\)`)
)

func New(pathProvider paths.Pather, baseRoute route.Route, files []*model.File) *FilePreviewExtension {
	return &FilePreviewExtension{
		pathProvider: pathProvider,
		base:         baseRoute,
		files:        files,
	}
}

type FilePreviewExtension struct {
	pathProvider paths.Pather
	base         route.Route
	files        []*model.File
}

func (converter *FilePreviewExtension) Convert(markdown string) (convertedContent string, converterError error) {

	convertedContent = markdown

	for {

		found, matches := pattern.IsMatch(convertedContent, markdownPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get the code
		renderedCode := converter.getGalleryCode(title, path)

		// replace markdown
		convertedContent = strings.Replace(convertedContent, originalText, renderedCode, 1)
	}

	return convertedContent, nil
}

func (converter *FilePreviewExtension) getGalleryCode(galleryTitle, path string) string {

	imageLinks := converter.getImageLinksByPath(galleryTitle, path)
	return fmt.Sprintf(`<section class="imagegallery">
				<header>%s</header>
				<ol>
					<li>
					%s
					</li>
				</ol>
			</section>`, galleryTitle, strings.Join(imageLinks, "\n</li>\n<li>\n"))
}

func (converter *FilePreviewExtension) getImageLinksByPath(galleryTitle, path string) []string {

	galleryRoute, err := route.NewFromRequest(path)
	if err != nil {
		// todo: log error
		return []string{}
	}

	baseRoute := converter.base
	fullGalleryRoute, err := route.Combine(baseRoute, galleryRoute)
	if err != nil {
		// todo: log error
		return []string{}
	}

	numberOfFiles := len(converter.files)
	imagelinks := make([]string, numberOfFiles, numberOfFiles)

	for index, file := range converter.files {

		// skip files which are not a child of the supplied path
		if !file.Route().IsChildOf(fullGalleryRoute) {
			continue
		}

		// skip files which are not images
		if !util.IsImageFile(file) {
			continue
		}

		imagePath := converter.pathProvider.Path(file.Route().Value())
		imageTitle := fmt.Sprintf("%s - %s (Image %v of %v)", galleryTitle, file.Route().LastComponentName(), index+1, numberOfFiles)

		imagelinks[index] = fmt.Sprintf(`<a href="%s" title="%s"><img src="%s" /></a>`, imagePath, imageTitle, imagePath)
	}

	return imagelinks
}
