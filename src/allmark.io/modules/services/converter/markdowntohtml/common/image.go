// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/route"
	"allmark.io/modules/services/thumbnail"
	"fmt"
	"strings"
)

func NewImageProvider(pathProvider paths.Pather, thumbnailIndex *thumbnail.Index) *ImageProvider {
	return &ImageProvider{
		pathProvider:   pathProvider,
		thumbnailIndex: thumbnailIndex,
	}
}

type ImageProvider struct {
	pathProvider   paths.Pather
	thumbnailIndex *thumbnail.Index
}

// GetImagePath returns the image path for the given file route.
// If one or more thumbnais exist it will return the thumbnail path (e.g. srcset="/thumbnails/105-D6134C1B-320-240.png 320w, /thumbnails/105-D6134C1B-640-480.png 640w, /thumbnails/105-D6134C1B-1024-768.png 1024w").
// If there is no thumbnail is will just return the canonical image path (e.g. src="document/files/sample.png")
func (provider *ImageProvider) GetImagePath(fileRoute route.Route) string {

	fullSizeImagePath := provider.getImagePath(fileRoute)

	// get thumbnail paths
	small, smallExists := provider.getThumbnailPath(fileRoute, thumbnail.SizeSmall)
	medium, mediumExists := provider.getThumbnailPath(fileRoute, thumbnail.SizeMedium)
	large, largeExists := provider.getThumbnailPath(fileRoute, thumbnail.SizeLarge)

	// assemble the image code
	imagePath := ""

	// assemble the src sets
	if smallExists || mediumExists || largeExists {

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

		if len(srcSets) > 0 {
			imagePath += fmt.Sprintf(` srcset="%s"`, strings.Join(srcSets, `, `))
		}
	}

	// use the full image as the defaults
	imagePath += fmt.Sprintf(` src="%s"`, fullSizeImagePath)

	return imagePath
}

func (provider *ImageProvider) getThumbnailPath(fileRoute route.Route, dimensions thumbnail.ThumbDimension) (thumbnailPath string, thumbnailAvailable bool) {

	// check if there are thumbs for the supplied file route
	thumbs, exists := provider.thumbnailIndex.GetThumbs(fileRoute.Value())
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
