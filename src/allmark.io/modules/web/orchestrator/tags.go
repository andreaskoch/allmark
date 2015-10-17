// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/view/viewmodel"
	"fmt"
	"math"
	"net/url"
)

var (

	// the maximum number of cloud entry levels
	tagCloudEntryLevels = 6
)

type TagsOrchestrator struct {
	*Orchestrator

	// caches and indizes
	tags     []viewmodel.Tag
	tagCloud viewmodel.TagCloud
}

// GetTags returns a list of all known tag models.
func (orchestrator *TagsOrchestrator) GetTags() []viewmodel.Tag {

	if orchestrator.tags != nil {
		return orchestrator.tags
	}

	// updateTags creates a tags list and assigns it to the orchestrator cache.
	updateTags := func(route route.Route) {

		rootItem := orchestrator.rootItem()
		if rootItem == nil {
			orchestrator.logger.Fatal("No root item found")
		}

		// items by tag
		itemsByTag := make(map[string][]*viewmodel.Model)
		for _, item := range orchestrator.getAllItems() {

			itemViewModel := &viewmodel.Model{
				Base: getBaseModel(rootItem, item, orchestrator.relativePather(rootItem.Route()), orchestrator.config),
			}

			for _, tag := range item.MetaData.Tags {
				if items, exists := itemsByTag[tag]; exists {
					itemsByTag[tag] = append(items, itemViewModel)
				} else {
					itemsByTag[tag] = []*viewmodel.Model{itemViewModel}
				}
			}

		}

		// create tag models
		tags := make([]viewmodel.Tag, 0)
		for tag, items := range itemsByTag {

			// create view model
			tagModel := viewmodel.Tag{
				Name:   tag,
				Anchor: url.QueryEscape(tag),
				Route:  orchestrator.tagPather().Path(url.QueryEscape(tag)),
				Children: items,
			}

			// append to list
			tags = append(tags, tagModel)
		}

		// sort the tags
		viewmodel.SortTagBy(tagsByName).Sort(tags)

		orchestrator.tags = tags
	}

	asyncUpdate := func(route route.Route) {
		go updateTags(route)
	}

	// register update callbacks
	orchestrator.registerUpdateCallback("update tags", UpdateTypeNew, asyncUpdate)
	orchestrator.registerUpdateCallback("update tags", UpdateTypeModified, asyncUpdate)
	orchestrator.registerUpdateCallback("update tags", UpdateTypeDeleted, asyncUpdate)

	// build the cache
	updateTags(route.New())

	return orchestrator.tags
}

// GetTagCloud returns the latest tag cloud viewmodel.
func (orchestrator *TagsOrchestrator) GetTagCloud() viewmodel.TagCloud {

	if orchestrator.tagCloud != nil {
		return orchestrator.tagCloud
	}

	// updateTagCloud creates a new tag cloud and assigns it to the orchestrator cache.
	updateTagCloud := func(route route.Route) {
		cloud := make(viewmodel.TagCloud, 0)

		minNumberOfItems := 1
		maxNumberOfItems := 1

		for _, tag := range orchestrator.GetTags() {

			// calculate the number of items per tag
			numberItemsPerTag := len(tag.Children)

			// update the maximum number of items per tag
			if numberItemsPerTag > maxNumberOfItems {
				maxNumberOfItems = numberItemsPerTag
			}

			// update the minimum number of items per tag
			if numberItemsPerTag < minNumberOfItems {
				minNumberOfItems = numberItemsPerTag
			}

			// create a new tag cloud entry
			tagCloudEntry := viewmodel.TagCloudEntry{
				Name:           tag.Name,
				Anchor:         url.QueryEscape(tag.Name),
				Route:          orchestrator.tagPather().Path(url.QueryEscape(tag.Name)),
				NumberOfChildren: numberItemsPerTag,
			}

			cloud = append(cloud, tagCloudEntry)
		}

		// update the tag cloud entry levels according
		// to the recorded min and max number of items
		for index, entry := range cloud {
			// calculate the entry level
			cloud[index].Level = getTagCloudEntryLevel(entry.NumberOfChildren, minNumberOfItems, maxNumberOfItems, tagCloudEntryLevels)
		}

		// sort tags by name
		viewmodel.SortTagCloudBy(tagCloudEntriesByName).Sort(cloud)

		orchestrator.tagCloud = cloud
	}

	asyncUpdate := func(route route.Route) {
		go updateTagCloud(route)
	}

	// register update callbacks
	orchestrator.registerUpdateCallback("update tagcloud", UpdateTypeNew, asyncUpdate)
	orchestrator.registerUpdateCallback("update tagcloud", UpdateTypeModified, asyncUpdate)
	orchestrator.registerUpdateCallback("update tagcloud", UpdateTypeDeleted, asyncUpdate)

	// build the cache
	updateTagCloud(route.New())

	return orchestrator.tagCloud
}

func (orchestrator *TagsOrchestrator) getItemTags(route route.Route) []viewmodel.Tag {

	var tags []viewmodel.Tag

	// abort if the item has no tags
	item := orchestrator.getItem(route)
	if item == nil {
		return tags
	}

	for _, tag := range item.MetaData.Tags {

		// create view model
		tagModel := viewmodel.Tag{
			Name:   tag,
			Anchor: url.QueryEscape(tag),
			Route:  orchestrator.tagPather().Path(url.QueryEscape(tag)),
		}

		// append to list
		tags = append(tags, tagModel)
	}

	return tags
}

func getTagCloudEntryLevel(numberOfChildren, minNumberOfChildren, maxNumberOfChildren, levelCount int) int {

	// check the number of children for negative numbers
	if numberOfChildren < 1 {
		panic(fmt.Sprintf("The number of children '%v' cannot be less than 1.", numberOfChildren))
	}

	// check max boundary
	if numberOfChildren > maxNumberOfChildren {
		panic(fmt.Sprintf("The number of children '%v' cannot be greater than the maximum number of children '%v'.", numberOfChildren, maxNumberOfChildren))
	}

	// check min boundary
	if numberOfChildren < minNumberOfChildren {
		panic(fmt.Sprintf("The number of children '%v' cannot be smaller than the minimum number of children '%v'.", numberOfChildren, minNumberOfChildren))
	}

	// check the level count
	if levelCount < 0 {
		panic(fmt.Sprintf("The level count must be greater than 0.", levelCount))
	}

	// calculate the ratio between the "number of children" to the "maximum number" of children (0, 1]
	ratioNumberOfChildrenToMaxNumberOfChildren := float64(numberOfChildren) / float64(maxNumberOfChildren)
	inverseRation := 1 - ratioNumberOfChildrenToMaxNumberOfChildren

	// calculate the level
	level := int(math.Floor(inverseRation*float64(levelCount))) + 1

	return level
}

// sort tags by nametag
func tagsByName(tag1, tag2 viewmodel.Tag) bool {
	return tag1.Name < tag2.Name
}

// sort tag cloud entries by name
func tagCloudEntriesByName(tagCloudEntry1, tagCloudEntry2 viewmodel.TagCloudEntry) bool {
	return tagCloudEntry1.Name < tagCloudEntry2.Name
}
