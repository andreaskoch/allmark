// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

func MapSitemap(item *repository.Item) *view.Sitemap {

	// map the childs
	childs := make([]*view.Sitemap, 0)
	for _, child := range item.Childs {
		childs = append(childs, MapSitemap(child))
	}

	// map the item
	return &view.Sitemap{
		AbsoluteRoute: item.AbsoluteRoute,
		Title:         item.Title,
		Description:   item.Description,
		Childs:        childs,
	}

}
