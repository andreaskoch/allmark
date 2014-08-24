// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
)

type XmlSitemapOrchestrator struct {
	*Orchestrator
}

func (orchestrator *XmlSitemapOrchestrator) GetSitemapEntires(hostname string) []viewmodel.XmlSitemapEntry {

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found")
	}

	childs := make([]viewmodel.XmlSitemapEntry, 0)
	for _, child := range orchestrator.repository.Items() {

		parsedItem := orchestrator.parseItem(child)
		if parsedItem == nil {
			orchestrator.logger.Warn("Cannot parse item %q", child.String())
			continue
		}

		// skip virtual items
		if parsedItem.IsVirtual() {
			continue
		}

		// item location
		addressPrefix := fmt.Sprintf("http://%s/", hostname)
		pathProvider := orchestrator.absolutePather(addressPrefix)
		location := pathProvider.Path(parsedItem.Route().Value())

		// last modified date
		lastModifiedDate := ""
		if parsedItem.MetaData != nil && parsedItem.MetaData.LastModifiedDate != nil {
			lastModifiedDate = parsedItem.MetaData.LastModifiedDate.Format("2006-01-02")
		}

		childs = append(childs, viewmodel.XmlSitemapEntry{
			Loc:          location,
			LastModified: lastModifiedDate,
		})
	}

	return childs
}
