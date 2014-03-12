// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func NewXmlSitemapOrchestrator(itemIndex *index.ItemIndex) XmlSitemapOrchestrator {
	return XmlSitemapOrchestrator{
		itemIndex: itemIndex,
	}
}

type XmlSitemapOrchestrator struct {
	itemIndex *index.ItemIndex
}

func (orchestrator *XmlSitemapOrchestrator) GetSitemapEntires(pathProvider paths.Pather) []viewmodel.XmlSitemapEntry {

	rootItem := orchestrator.itemIndex.Root()
	if rootItem == nil {
		panic("No root item found")
	}

	childs := make([]viewmodel.XmlSitemapEntry, 0)
	for _, child := range orchestrator.itemIndex.GetAllChilds(rootItem.Route()) {

		childs = append(childs, viewmodel.XmlSitemapEntry{
			Loc: child.Route().Value(),
		})
	}

	return childs
}
