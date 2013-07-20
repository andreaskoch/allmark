// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

func MapSitemap(root *repository.Item) *view.Sitemap {
	return mapSitemap(root)
}

func mapSitemap(item *repository.Item) *view.Sitemap {

	// map the childs
	childs := make([]*view.Sitemap, 0)
	for _, child := range item.Childs {
		childs = append(childs, mapSitemap(child))
	}

	// map the item
	return &view.Sitemap{
		AbsoluteRoute: item.AbsoluteRoute,
		RelativeRoute: item.RelativeRoute,
		Title:         item.Title,
		Description:   item.Description,
		Childs:        childs,
	}

}
