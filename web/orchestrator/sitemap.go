// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
)

type SitemapOrchestrator struct {
	*Orchestrator

	// caches
	sitemap *viewmodel.Sitemap
}

func (orchestrator *SitemapOrchestrator) GetSitemap() viewmodel.Sitemap {

	cacheType := "html sitmap"

	// load from cache
	if orchestrator.sitemap != nil {

		// re-prime the cache if it is stale
		if orchestrator.isCacheStale(cacheType) {
			go orchestrator.primeCache(cacheType)
		}

		return *orchestrator.sitemap
	}

	orchestrator.setCache(cacheType, func() {

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
	})

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
