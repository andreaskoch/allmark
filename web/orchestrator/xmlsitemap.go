// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"time"
)

type XmlSitemapOrchestrator struct {
	*Orchestrator
}

func (orchestrator *XmlSitemapOrchestrator) GetSitemapEntires(hostname string) []viewmodel.XmlSitemapEntry {

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found")
	}

	zeroTime := time.Time{}

	childs := make([]viewmodel.XmlSitemapEntry, 0)
	for _, item := range orchestrator.getAllItems() {

		// skip virtual items
		if item.IsVirtual() {
			continue
		}

		// item location
		addressPrefix := fmt.Sprintf("http://%s/", hostname)
		pathProvider := orchestrator.absolutePather(addressPrefix)
		location := pathProvider.Path(item.Route().Value())

		// last modified date
		lastModifiedDate := ""
		if item.MetaData.LastModifiedDate != zeroTime {
			lastModifiedDate = item.MetaData.LastModifiedDate.Format("2006-01-02")
		}

		childs = append(childs, viewmodel.XmlSitemapEntry{
			Loc:          location,
			LastModified: lastModifiedDate,
		})
	}

	return childs
}
