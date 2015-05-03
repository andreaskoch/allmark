// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imagegallery

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
	// imagegallery: [*description text*](*folder path*)
	markdownPattern = regexp.MustCompile(`imagegallery: \[([^\]]*)\]\(([^)]+)\)`)
)

func New(pathProvider paths.Pather, baseRoute route.Route, files []*model.File, imageProvider *common.ImageProvider) *FilePreviewExtension {
	return &FilePreviewExtension{
		pathProvider:  pathProvider,
		base:          baseRoute,
		files:         files,
		imageProvider: imageProvider,
	}
}

type FilePreviewExtension struct {
	pathProvider  paths.Pather
	base          route.Route
	files         []*model.File
	imageProvider *common.ImageProvider
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

	imageLinks := converter.getImageLinksByPath(path)
	if galleryTitle != "" {
		return fmt.Sprintf(`<section class="imagegallery">
					<header>%s</header>
					<ol>
						<li>
						%s
						</li>
					</ol>
				</section>`, galleryTitle, strings.Join(imageLinks, "\n</li>\n<li>\n"))
	}

	return fmt.Sprintf(`<section class="imagegallery">
					<ol>
						<li>
						%s
						</li>
					</ol>
				</section>`, strings.Join(imageLinks, "\n</li>\n<li>\n"))
}

func (converter *FilePreviewExtension) getImageLinksByPath(path string) []string {

	baseRoute := converter.base
	galleryRoute := route.NewFromRequest(path)
	fullGalleryRoute := route.Combine(baseRoute, galleryRoute)

	numberOfFiles := len(converter.files)
	imagelinks := make([]string, 0, numberOfFiles)

	for _, file := range converter.files {

		// skip files which are not a child of the supplied path
		if !file.Route().IsChildOf(fullGalleryRoute) {
			continue
		}

		// skip files which are not images
		if !util.IsImageFile(file) {
			continue
		}

		// image title
		imageTitle := file.Route().LastComponentName() // use the file name for the title

		// calculate the image code
		imageCode := converter.imageProvider.GetImageCodeWithLink(imageTitle, file.Route())
		imagelinks = append(imagelinks, imageCode)
	}

	return imagelinks
}
