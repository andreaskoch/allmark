// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func NewXmlSitemapOrchestrator(itemIndex *index.Index) XmlSitemapOrchestrator {
	return XmlSitemapOrchestrator{
		itemIndex: itemIndex,
	}
}

type XmlSitemapOrchestrator struct {
	itemIndex *index.Index
}

func (orchestrator *XmlSitemapOrchestrator) GetSitemapEntires(pathProvider paths.Pather) []viewmodel.XmlSitemapEntry {

	rootItem := orchestrator.itemIndex.Root()
	if rootItem == nil {
		panic("No root item found")
	}

	childs := make([]viewmodel.XmlSitemapEntry, 0)
	for _, child := range orchestrator.itemIndex.Items() {

		// skip virtual items
		if child.IsVirtual() {
			continue
		}

		// item location
		location := pathProvider.Path(child.Route().Value())

		// last modified date
		lastModifiedDate := ""
		if child.MetaData != nil && child.MetaData.LastModifiedDate != nil {
			lastModifiedDate = child.MetaData.LastModifiedDate.Format("2006-01-02")
		}

		childs = append(childs, viewmodel.XmlSitemapEntry{
			Loc:          location,
			LastModified: lastModifiedDate,
		})
	}

	return childs
}
