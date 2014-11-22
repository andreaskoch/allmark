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

func (orchestrator *FeedOrchestrator) GetEntries(hostname string, itemsPerPage, page int) (entries []viewmodel.FeedEntry, found bool) {

	// validate page number
	if page < 1 {
		orchestrator.logger.Fatal("Invalid page number (%v).", page)
	}

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found")
	}

	// create the path provider
	addressPrefix := fmt.Sprintf("http://%s/", hostname)
	pathProvider := orchestrator.absolutePather(addressPrefix)

	feedEntries := make([]viewmodel.FeedEntry, 0)

	latestItems, found := pagedItems(orchestrator.getLatestItems(rootItem.Route()), itemsPerPage, page)
	if !found {
		return feedEntries, false
	}

	for _, item := range latestItems {
		feedEntries = append(feedEntries, orchestrator.createFeedEntryModel(pathProvider, item))
	}

	return feedEntries, true
}

func (orchestrator *FeedOrchestrator) createFeedEntryModel(pathProvider paths.Pather, item *model.Item) viewmodel.FeedEntry {

	// item location
	location := pathProvider.Path(item.Route().Value())

	// content
	content, err := orchestrator.converter.Convert(orchestrator.getItemByAlias, pathProvider, item)
	if err != nil {
		content = err.Error()
	}

	// creation date
	creationDate := item.MetaData.CreationDate.Format("2006-01-02")

	return viewmodel.FeedEntry{
		Title:       item.Title,
		Description: content,
		Link:        location,
		PubDate:     creationDate,
	}
}
