// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package preprocessor

import (
	"github.com/andreaskoch/allmark/common/paths"
	"github.com/andreaskoch/allmark/common/route"
	"github.com/andreaskoch/allmark/model"
	"github.com/andreaskoch/allmark/services/converter/markdowntohtml/imageprovider"
	"fmt"
	"regexp"
	"strings"
)

var (
	// imagegallery: [*description text*](*folder path*)
	imageGalleryExtensionPattern = regexp.MustCompile(`imagegallery: \[([^\]]*)\]\(([^)]+)\)`)
)

func newImageGalleryExtension(pathProvider paths.Pather, baseRoute route.Route, files []*model.File, imageProvider *imageprovider.ImageProvider) *imageGalleryExtension {
	return &imageGalleryExtension{
		pathProvider:  pathProvider,
		base:          baseRoute,
		files:         files,
		imageProvider: imageProvider,
	}
}

type imageGalleryExtension struct {
	pathProvider  paths.Pather
	base          route.Route
	files         []*model.File
	imageProvider *imageprovider.ImageProvider
}

func (converter *imageGalleryExtension) Convert(markdown string) (convertedContent string, converterError error) {

	convertedContent = markdown

	for _, match := range imageGalleryExtensionPattern.FindAllStringSubmatch(convertedContent, -1) {

		if len(match) != 3 {
			continue
		}

		// parameters
		originalText := strings.TrimSpace(match[0])
		title := strings.TrimSpace(match[1])
		path := strings.TrimSpace(match[2])

		// get the code
		renderedCode := converter.getGalleryCode(title, path)

		// replace markdown
		convertedContent = strings.Replace(convertedContent, originalText, renderedCode, 1)
	}

	return convertedContent, nil
}

func (converter *imageGalleryExtension) getGalleryCode(galleryTitle, path string) string {

	imageLinks := converter.getImageLinksByPath(path)

	var code string
	if galleryTitle != "" {
		code += fmt.Sprintf("**%s**\n\n", galleryTitle)
	}
	code += strings.Join(imageLinks, "\n")

	return code
}

func (converter *imageGalleryExtension) getImageLinksByPath(path string) []string {

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
		if !model.IsImageFile(file) {
			continue
		}

		// image title
		imageTitle := file.Route().LastComponentName() // use the file name for the title

		// calculate the image code
		imagePath := converter.imageProvider.GetImagePath(converter.pathProvider, file.Route())
		imageCode := fmt.Sprintf(`<img %s alt="%s"/>`, imagePath, imageTitle)

		// link the image to the original
		fullSizeImagePath := converter.pathProvider.Path(file.Route().Value())
		imageWithLink := fmt.Sprintf(`<a href="%s" title="%s">%s</a>`, fullSizeImagePath, imageTitle, imageCode)

		imagelinks = append(imagelinks, imageWithLink)
	}

	return imagelinks
}
