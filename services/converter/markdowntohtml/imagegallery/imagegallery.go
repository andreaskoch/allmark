// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imagegallery

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/pattern"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/converter/markdowntohtml/util"
	"github.com/andreaskoch/allmark2/services/thumbnail"
	"regexp"
	"strings"
)

var (
	// imagegallery: [*description text*](*folder path*)
	markdownPattern = regexp.MustCompile(`imagegallery: \[([^\]]*)\]\(([^)]+)\)`)
)

func New(pathProvider paths.Pather, baseRoute route.Route, files []*model.File, thumbnailIndex *thumbnail.Index) *FilePreviewExtension {
	return &FilePreviewExtension{
		pathProvider:   pathProvider,
		base:           baseRoute,
		files:          files,
		thumbnailIndex: thumbnailIndex,
	}
}

type FilePreviewExtension struct {
	pathProvider   paths.Pather
	base           route.Route
	files          []*model.File
	thumbnailIndex *thumbnail.Index
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

		// get paths
		fullSizeImagePath := converter.getImagePath(file.Route())
		thumnailPath := converter.getThumbnailPath(file.Route())

		// image title
		imageTitle := file.Route().LastComponentName() // file name

		imagelinks[index] = fmt.Sprintf(`<a href="%s" title="%s"><img src="%s" /></a>`, fullSizeImagePath, imageTitle, thumnailPath)
	}

	return imagelinks
}

func (converter *FilePreviewExtension) getThumbnailPath(fileRoute route.Route) string {

	// assemble to the full image route
	fullRoute, err := route.Combine(converter.base, fileRoute)
	if err != nil {
		panic(fmt.Sprintf("Cannot combine routes %q and %q.", converter.base, fileRoute))
	}

	// check if there are thumbs for the supplied file route
	thumbs, exists := converter.thumbnailIndex.GetThumbs(fullRoute.Value())
	if !exists {
		return converter.getImagePath(fileRoute) // return the full-size image path
	}

	// lookup thumb by size
	thumb, exists := thumbs.GetThumbBySize(thumbnail.SizeMedium)
	if !exists {
		return converter.getImagePath(fileRoute) // return the full-size image path

	}

	return converter.getImagePath(thumb.ThumbRoute())
}

func (converter *FilePreviewExtension) getImagePath(fileRoute route.Route) string {
	return converter.pathProvider.Path(fileRoute.Value())
}
