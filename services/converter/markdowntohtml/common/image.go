// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/thumbnail"
	"strings"
)

func NewImageProvider(pathProvider paths.Pather, baseRoute route.Route, files []*model.File, thumbnailIndex *thumbnail.Index) *ImageProvider {
	return &ImageProvider{
		pathProvider:   pathProvider,
		base:           baseRoute,
		thumbnailIndex: thumbnailIndex,
	}
}

type ImageProvider struct {
	pathProvider   paths.Pather
	base           route.Route
	thumbnailIndex *thumbnail.Index
}

func (provider *ImageProvider) GetImageCodeWithLink(imageTitle string, fileRoute route.Route) string {
	fullSizeImagePath := provider.getImagePath(fileRoute)
	imageCode := provider.GetImageCode(imageTitle, fileRoute)
	return fmt.Sprintf(`<a href="%s" title="%s">%s</a>`, fullSizeImagePath, imageTitle, imageCode)
}

func (provider *ImageProvider) GetImageCode(imageTitle string, fileRoute route.Route) string {

	fullSizeImagePath := provider.getImagePath(fileRoute)

	// get thumbnail paths
	small, smallExists := provider.getThumbnailPath(fileRoute, thumbnail.SizeSmall)
	medium, mediumExists := provider.getThumbnailPath(fileRoute, thumbnail.SizeMedium)
	large, largeExists := provider.getThumbnailPath(fileRoute, thumbnail.SizeLarge)

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

func (provider *ImageProvider) getThumbnailPath(fileRoute route.Route, dimensions thumbnail.ThumbDimension) (thumbnailPath string, thumbnailAvailable bool) {

	// assemble to the full image route
	fullRoute, err := route.Combine(provider.base, fileRoute)
	if err != nil {
		panic(fmt.Sprintf("Cannot combine routes %q and %q.", provider.base, fileRoute))
	}

	// check if there are thumbs for the supplied file route
	thumbs, exists := provider.thumbnailIndex.GetThumbs(fullRoute.Value())
	if !exists {
		return "", false // return the full-size image path
	}

	// lookup thumb by size
	thumb, exists := thumbs.GetThumbBySize(dimensions)
	if !exists {
		return "", false // return the full-size image path

	}

	return provider.getImagePath(thumb.ThumbRoute()), true
}

func (provider *ImageProvider) getImagePath(fileRoute route.Route) string {
	return provider.pathProvider.Path(fileRoute.Value())
}
