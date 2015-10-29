// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"allmark.io/modules/model"
	"allmark.io/modules/web/view/viewmodel"
	"fmt"
)

// A FeedOrchestrator provides feed models.
type FeedOrchestrator struct {
	*Orchestrator
}

// GetFeed returns a feed model for the given base URL, items per page and page.
func (orchestrator *FeedOrchestrator) GetFeed(baseURL string, itemsPerPage, page int) (viewmodel.Feed, error) {
	root, err := orchestrator.getRootEntry(baseURL)
	if err != nil {
		return viewmodel.Feed{}, err
	}

	items, err := orchestrator.getItems(baseURL, itemsPerPage, page)
	if err != nil {
		return viewmodel.Feed{}, err
	}

	feedModel := viewmodel.Feed{}
	feedModel.FeedEntry = root
	feedModel.Items = items

	return feedModel, nil
}

func (orchestrator *FeedOrchestrator) getRootEntry(baseURL string) (viewmodel.FeedEntry, error) {

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		return viewmodel.FeedEntry{}, fmt.Errorf("No root item found.")
	}

	return orchestrator.createFeedEntryModel(baseURL, rootItem), nil
}

func (orchestrator *FeedOrchestrator) getItems(baseURL string, itemsPerPage, page int) ([]viewmodel.FeedEntry, error) {

	// validate page number
	if page < 1 {
		return []viewmodel.FeedEntry{}, fmt.Errorf("Invalid page number: %v.", page)
	}

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		return []viewmodel.FeedEntry{}, fmt.Errorf("No root item found.")
	}

	var feedEntries []viewmodel.FeedEntry

	latestItems, found := pagedItems(orchestrator.getLatestItems(rootItem.Route()), itemsPerPage, page)
	if !found {
		return feedEntries, fmt.Errorf("No items found (Items per page: %v, Page: %v)", itemsPerPage, page)
	}

	for _, item := range latestItems {
		feedEntries = append(feedEntries, orchestrator.createFeedEntryModel(baseURL, item))
	}

	return feedEntries, nil
}

func (orchestrator *FeedOrchestrator) createFeedEntryModel(baseURL string, item *model.Item) viewmodel.FeedEntry {

	rootPathProvider := orchestrator.absolutePather(fmt.Sprintf("%s/", baseURL))

	// item location
	location := rootPathProvider.Path(item.Route().Value())

	// content
	content, err := orchestrator.converter.Convert(orchestrator.getItemByAlias, rootPathProvider, item)
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
