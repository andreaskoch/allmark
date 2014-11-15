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

		// image title
		imageTitle := file.Route().LastComponentName() // use the file name for the title

		// get paths
		fullSizeImagePath := converter.getImagePath(file.Route())

		// get the image code
		imageCode := converter.getImageCode(imageTitle, fullSizeImagePath, file.Route())

		imagelinks[index] = fmt.Sprintf(`<a href="%s" title="%s">%s</a>`, fullSizeImagePath, imageTitle, imageCode)
	}

	return imagelinks
}

func (converter *FilePreviewExtension) getImageCode(imageTitle, fullSizeImagePath string, fileRoute route.Route) string {

	// get thumbnail paths
	small, smallExists := converter.getThumbnailPath(fileRoute, thumbnail.SizeSmall)
	medium, mediumExists := converter.getThumbnailPath(fileRoute, thumbnail.SizeMedium)
	large, largeExists := converter.getThumbnailPath(fileRoute, thumbnail.SizeLarge)

	// assemble the image code
	image := "<img"

	// assemble the src sets
	if smallExists || mediumExists || largeExists {

		image += " srcset=\""

		srcSets := make([]string, 0)
		if smallExists {
			srcSets = append(srcSets, small+fmt.Sprintf(" %vw", thumbnail.SizeSmall.MaxWidth))
		}

		if mediumExists {
			srcSets = append(srcSets, medium+fmt.Sprintf(" %vw", thumbnail.SizeMedium.MaxWidth))
		}

		if largeExists {
			srcSets = append(srcSets, large+fmt.Sprintf(" %vw", thumbnail.SizeLarge.MaxWidth))
		}

		image += strings.Join(srcSets, ", ")
	}

	// default image
	if smallExists || mediumExists || largeExists {

		// use the small image as the default
		image += " src=\"" + small + "\""

	} else {

		// use the full image as the defaults
		image += " src=\"" + fullSizeImagePath + "\""

	}

	image += fmt.Sprintf(` alt="%s" />`, imageTitle)

	return image
}

func (converter *FilePreviewExtension) getThumbnailPath(fileRoute route.Route, dimensions thumbnail.ThumbDimension) (thumbnailPath string, thumbnailAvailable bool) {

	// assemble to the full image route
	fullRoute, err := route.Combine(converter.base, fileRoute)
	if err != nil {
		panic(fmt.Sprintf("Cannot combine routes %q and %q.", converter.base, fileRoute))
	}

	// check if there are thumbs for the supplied file route
	thumbs, exists := converter.thumbnailIndex.GetThumbs(fullRoute.Value())
	if !exists {
		return "", false // return the full-size image path
	}

	// lookup thumb by size
	thumb, exists := thumbs.GetThumbBySize(dimensions)
	if !exists {
		return "", false // return the full-size image path

	}

	return converter.getImagePath(thumb.ThumbRoute()), true
}

func (converter *FilePreviewExtension) getImagePath(fileRoute route.Route) string {
	return converter.pathProvider.Path(fileRoute.Value())
}
