// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
)

type FeedOrchestrator struct {
	*Orchestrator
}

func (orchestrator *FeedOrchestrator) GetRootEntry(hostname string) viewmodel.FeedEntry {

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found.")
	}

	addressPrefix := fmt.Sprintf("http://%s/", hostname)
	pathProvider := orchestrator.absolutePather(addressPrefix)

	return orchestrator.createFeedEntryModel(pathProvider, rootItem)
}

func (orchestrator *FeedOrchestrator) GetEntries(hostname string, itemsPerPage, page int) []viewmodel.FeedEntry {

	// validate page number
	if page < 1 {
		orchestrator.logger.Fatal("Invalid page number (%v).", page)
	}

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found")
	}

	// determine start item
	startItemNumber := itemsPerPage * (page - 1)

	// determine end item
	endItemNumber := itemsPerPage * page

	// create the path provider
	addressPrefix := fmt.Sprintf("http://%s/", hostname)
	pathProvider := orchestrator.absolutePather(addressPrefix)

	childs := make([]viewmodel.FeedEntry, 0)
	for _, item := range orchestrator.repository.Items() {

		parsedItem := orchestrator.parseItem(item)
		if parsedItem == nil {
			orchestrator.logger.Warn("Cannot parse item %q", item.String())
			continue
		}

		// skip virtual items
		if parsedItem.IsVirtual() {
			continue
		}

		// paging
		currentNumberOfItems := len(childs)
		if currentNumberOfItems < startItemNumber || currentNumberOfItems >= endItemNumber {
			continue
		}

		childs = append(childs, orchestrator.createFeedEntryModel(pathProvider, parsedItem))
	}

	return childs
}

func (orchestrator *FeedOrchestrator) createFeedEntryModel(pathProvider paths.Pather, item *model.Item) viewmodel.FeedEntry {

	// item location
	location := pathProvider.Path(item.Route().Value())

	// content
	content, err := orchestrator.converter.Convert(pathProvider, item)
	if err != nil {
		content = err.Error()
	}

	// last modified date
	lastModifiedDate := item.MetaData.LastModifiedDate.Format("2006-01-02")

	return viewmodel.FeedEntry{
		Title:       item.Title,
		Description: content,
		Link:        location,
		PubDate:     lastModifiedDate,
	}
}
