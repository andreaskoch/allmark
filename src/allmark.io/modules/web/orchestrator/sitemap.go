// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/view/viewmodel"
)

type SitemapOrchestrator struct {
	*Orchestrator

	// caches
	sitemap *viewmodel.Sitemap
}

func (orchestrator *SitemapOrchestrator) GetSitemap() viewmodel.Sitemap {

	if orchestrator.sitemap != nil {
		return *orchestrator.sitemap
	}

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found")
	}

	sitemapModel := viewmodel.Sitemap{
		Title:       rootItem.Title,
		Description: rootItem.Description,
		Childs:      orchestrator.getSitemapEntries(rootItem.Route()),
		Path:        "/",
	}

	orchestrator.sitemap = &sitemapModel

	return *orchestrator.sitemap
}

func (orchestrator *SitemapOrchestrator) getSitemapEntries(startRoute route.Route) []viewmodel.Sitemap {

	childs := make([]viewmodel.Sitemap, 0)
	for _, child := range orchestrator.getChilds(startRoute) {

		childRoute := child.Route()

		childModel := viewmodel.Sitemap{
			Title:       child.Title,
			Description: child.Description,
			Childs:      orchestrator.getSitemapEntries(childRoute),
			Path:        orchestrator.itemPather().Path(childRoute.Value()),
		}

		childs = append(childs, childModel)
	}

	return childs
}
