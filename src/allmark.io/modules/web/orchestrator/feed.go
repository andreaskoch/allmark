// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"allmark.io/modules/model"
	"allmark.io/modules/web/view/viewmodel"
	"fmt"
)

type FeedOrchestrator struct {
	*Orchestrator
}

func (orchestrator *FeedOrchestrator) GetRootEntry(baseURL string) viewmodel.FeedEntry {

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found.")
	}

	return orchestrator.createFeedEntryModel(baseURL, rootItem)
}

func (orchestrator *FeedOrchestrator) GetEntries(baseURL string, itemsPerPage, page int) (entries []viewmodel.FeedEntry, found bool) {

	// validate page number
	if page < 1 {
		orchestrator.logger.Fatal("Invalid page number (%v).", page)
	}

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found")
	}

	feedEntries := make([]viewmodel.FeedEntry, 0)

	latestItems, found := pagedItems(orchestrator.getLatestItems(rootItem.Route()), itemsPerPage, page)
	if !found {
		return feedEntries, false
	}

	for _, item := range latestItems {
		feedEntries = append(feedEntries, orchestrator.createFeedEntryModel(baseURL, item))
	}

	return feedEntries, true
}

func (orchestrator *FeedOrchestrator) createFeedEntryModel(baseURL string, item *model.Item) viewmodel.FeedEntry {

	rootPathProvider := orchestrator.absolutePather(fmt.Sprintf("%s/", baseURL))
	itemContentPathProvider := orchestrator.absolutePather(fmt.Sprintf("%s/%s/", baseURL, item.Route().Value()))

	// item location
	location := rootPathProvider.Path(item.Route().Value())

	// content
	content, err := orchestrator.converter.Convert(orchestrator.getItemByAlias, rootPathProvider, itemContentPathProvider, item)
	if err != nil {
		content = err.Error()
	}

	// append the description
	if item.Description != "" {
		content = fmt.Sprintf("<p>%s</p>\n\n%s", item.Description, content)
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
