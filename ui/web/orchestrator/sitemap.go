// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	// "github.com/andreaskoch/allmark2/common/route"
	// "github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func NewSitemapOrchestrator(itemIndex *index.ItemIndex) SitemapOrchestrator {
	return SitemapOrchestrator{
		itemIndex: itemIndex,
	}
}

type SitemapOrchestrator struct {
	itemIndex *index.ItemIndex
}

func (orchestrator *SitemapOrchestrator) GetSitemap(pathProvider paths.Pather) viewmodel.Sitemap {

	// orchestrator.itemIndex.Walk(func(item *model.Item) {

	// 	viewModel := viewmodel.Sitemap{
	// 		Type:    item.Type.String(),
	// 		Route:   item.Route().Value(),
	// 		Level:   item.Route().Level(),
	// 		BaseUrl: getBaseUrlFromItem(item.Route()),

	// 		Title:       item.Title,
	// 		Description: item.Description,
	// 	}
	// })

	return viewmodel.Sitemap{}
}
