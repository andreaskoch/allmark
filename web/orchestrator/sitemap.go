// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/web/view/viewmodel"
)

type SitemapOrchestrator struct {
	*Orchestrator

	// caches
	sitemap *viewmodel.SitemapEntry
}

func (orchestrator *SitemapOrchestrator) GetSitemap() viewmodel.SitemapEntry {

	if orchestrator.sitemap != nil {
		return *orchestrator.sitemap
	}

	// updateSitemap creates a new sitemap model and assigns it to the orchestrator cache.
	updateSitemap := func(route route.Route) {
		rootItem := orchestrator.rootItem()
		if rootItem == nil {
			orchestrator.logger.Fatal("No root item found")
		}

		sitemapModel := viewmodel.SitemapEntry{
			Title:       rootItem.Title,
			Description: rootItem.Description,
			Children:      orchestrator.getSitemapEntries(rootItem.Route()),
			Path:        "/",
		}

		orchestrator.sitemap = &sitemapModel
	}

	// register update callbacks
	orchestrator.registerUpdateCallback("update sitemap", UpdateTypeNew, updateSitemap)
	orchestrator.registerUpdateCallback("update sitemap", UpdateTypeModified, updateSitemap)
	orchestrator.registerUpdateCallback("update sitemap", UpdateTypeDeleted, updateSitemap)

	// build the first sitemap
	updateSitemap(route.New())

	return *orchestrator.sitemap
}

func (orchestrator *SitemapOrchestrator) getSitemapEntries(startRoute route.Route) []viewmodel.SitemapEntry {

	children := make([]viewmodel.SitemapEntry, 0)
	for _, child := range orchestrator.getChildren(startRoute) {

		childRoute := child.Route()

		childModel := viewmodel.SitemapEntry{
			Title:       child.Title,
			Description: child.Description,
			Children:      orchestrator.getSitemapEntries(childRoute),
			Path:        orchestrator.itemPather().Path(childRoute.Value()),
		}

		children = append(children, childModel)
	}

	return children
}
