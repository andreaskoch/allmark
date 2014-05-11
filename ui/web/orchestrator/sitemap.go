// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func NewSitemapOrchestrator(itemIndex *index.Index) SitemapOrchestrator {
	return SitemapOrchestrator{
		itemIndex: itemIndex,
	}
}

type SitemapOrchestrator struct {
	itemIndex *index.Index
}

func (orchestrator *SitemapOrchestrator) GetSitemap(pathProvider paths.Pather) viewmodel.Sitemap {

	rootItem := orchestrator.itemIndex.Root()
	if rootItem == nil {
		panic("No root item found")
	}

	rootModel := viewmodel.Sitemap{
		Title:       rootItem.Title,
		Description: rootItem.Description,
		Childs:      getSitemapEntries(orchestrator.itemIndex, *rootItem.Route()),
		Path:        "/",
	}

	return rootModel
}

func getSitemapEntries(index *index.Index, startRoute route.Route) []viewmodel.Sitemap {

	childs := make([]viewmodel.Sitemap, 0)
	for _, child := range index.GetDirectChilds(&startRoute) {

		childRoute := child.Route()

		childModel := viewmodel.Sitemap{
			Title:       child.Title,
			Description: child.Description,
			Childs:      getSitemapEntries(index, *childRoute),
			Path:        child.Route().Value(),
		}

		childs = append(childs, childModel)
	}

	return childs
}
