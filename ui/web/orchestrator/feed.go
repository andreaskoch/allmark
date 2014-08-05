// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func NewFeedOrchestrator(itemIndex *index.Index, converter conversion.Converter) FeedOrchestrator {
	return FeedOrchestrator{
		itemIndex: itemIndex,
		converter: converter,
	}
}

type FeedOrchestrator struct {
	itemIndex *index.Index
	converter conversion.Converter
}

func (orchestrator *FeedOrchestrator) GetRootEntry(pathProvider paths.Pather) viewmodel.FeedEntry {

	rootItem := orchestrator.itemIndex.Root()
	if rootItem == nil {
		panic("No root item found")
	}

	return orchestrator.createFeedEntryModel(pathProvider, rootItem)
}

func (orchestrator *FeedOrchestrator) GetEntries(pathProvider paths.Pather, itemsPerPage, page int) []viewmodel.FeedEntry {

	// validate page number
	if page < 1 {
		panic("Invalid page number")
	}

	rootItem := orchestrator.itemIndex.Root()
	if rootItem == nil {
		panic("No root item found")
	}

	// determine start item
	startItemNumber := itemsPerPage * (page - 1)

	// determine end item
	endItemNumber := itemsPerPage * page

	childs := make([]viewmodel.FeedEntry, 0)
	for _, child := range orchestrator.itemIndex.Items() {

		// skip virtual items
		if child.IsVirtual() {
			continue
		}

		// paging
		currentNumberOfItems := len(childs)
		if currentNumberOfItems < startItemNumber || currentNumberOfItems >= endItemNumber {
			continue
		}

		childs = append(childs, orchestrator.createFeedEntryModel(pathProvider, child))
	}

	return childs
}

func (orchestrator *FeedOrchestrator) createFeedEntryModel(pathProvider paths.Pather, item *model.Item) viewmodel.FeedEntry {

	// item location
	location := pathProvider.Path(item.Route().Value())

	// content
	content, err := orchestrator.converter.Convert(pathProvider, item, false)
	if err != nil {
		content = err.Error()
	}

	// last modified date
	lastModifiedDate := ""
	if item.MetaData != nil && item.MetaData.LastModifiedDate != nil {
		lastModifiedDate = item.MetaData.LastModifiedDate.Format("2006-01-02")
	}

	return viewmodel.FeedEntry{
		Title:       item.Title,
		Description: content,
		Link:        location,
		PubDate:     lastModifiedDate,
	}
}
