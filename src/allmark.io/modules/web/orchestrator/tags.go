// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"allmark.io/modules/common/route"
	"allmark.io/modules/model"
	"allmark.io/modules/web/view/viewmodel"
	"fmt"
	"math"
)

var (

	// the maximum number of cloud entry levels
	tagCloudEntryLevels = 6
)

type TagsOrchestrator struct {
	*Orchestrator

	// caches and indizes
	tags     []*viewmodel.Tag
	tagCloud *viewmodel.TagCloud
}

func (orchestrator *TagsOrchestrator) GetTags() []*viewmodel.Tag {

	cacheType := "tags"

	// load from cache
	if orchestrator.tags != nil {

		// re-prime the cache if it is stale
		if orchestrator.isCacheStale(cacheType) {
			go orchestrator.primeCache(cacheType)
		}

		return orchestrator.tags
	}

	orchestrator.setCache(cacheType, func() {

		rootItem := orchestrator.rootItem()
		if rootItem == nil {
			orchestrator.logger.Fatal("No root item found")
		}

		// items by tag
		itemsByTag := make(map[string][]*viewmodel.Model)
		for _, item := range orchestrator.getAllItems() {

			itemViewModel := &viewmodel.Model{
				Base: getBaseModel(rootItem, item, orchestrator.relativePather(rootItem.Route())),
			}

			tags := []model.Tag{}
			if len(item.MetaData.Tags) > 0 {
				tags = item.MetaData.Tags
			}

			for _, tag := range tags {
				if items, exists := itemsByTag[tag.Name()]; exists {
					itemsByTag[tag.Name()] = append(items, itemViewModel)
				} else {
					itemsByTag[tag.Name()] = []*viewmodel.Model{itemViewModel}
				}
			}

		}

		// create tag models
		tags := make([]*viewmodel.Tag, 0)
		for tag, items := range itemsByTag {

			// create view model
			tagModel := &viewmodel.Tag{
				Name:   tag,
				Route:  orchestrator.tagPather().Path(tag),
				Childs: items,
			}

			// append to list
			tags = append(tags, tagModel)
		}

		// sort the tags
		viewmodel.SortTagBy(tagsByName).Sort(tags)

		orchestrator.tags = tags
	})

	return orchestrator.tags
}

func (orchestrator *TagsOrchestrator) GetTagCloud() *viewmodel.TagCloud {

	cacheType := "tagCloud"

	// load from cache
	if orchestrator.tagCloud != nil {

		// re-prime the cache if it is stale
		if orchestrator.isCacheStale(cacheType) {
			go orchestrator.primeCache(cacheType)
		}

		return orchestrator.tagCloud
	}

	orchestrator.setCache(cacheType, func() {

		cloud := make(viewmodel.TagCloud, 0)

		minNumberOfItems := 1
		maxNumberOfItems := 1

		for _, tag := range orchestrator.GetTags() {

			// calculate the number of items per tag
			numberItemsPerTag := len(tag.Childs)

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
				Route:          orchestrator.tagPather().Path(tag.Name),
				NumberOfChilds: numberItemsPerTag,
			}

			cloud = append(cloud, &tagCloudEntry)
		}

		// update the tag cloud entry levels according
		// to the recorded min and max number of items
		for index, entry := range cloud {
			// calculate the entry level
			cloud[index].Level = getTagCloudEntryLevel(entry.NumberOfChilds, minNumberOfItems, maxNumberOfItems, tagCloudEntryLevels)
		}

		// sort tags by name
		viewmodel.SortTagCloudBy(tagCloudEntriesByName).Sort(cloud)

		orchestrator.tagCloud = &cloud
	})

	return orchestrator.tagCloud
}

func (orchestrator *TagsOrchestrator) getItemTags(route route.Route) []*viewmodel.Tag {

	tags := make([]*viewmodel.Tag, 0)

	// abort if the item has no tags
	item := orchestrator.getItem(route)
	if item == nil {
		return tags
	}

	for _, tag := range item.MetaData.Tags {

		// create view model
		tagModel := &viewmodel.Tag{
			Name:  tag.Name(),
			Route: orchestrator.tagPather().Path(tag.Name()),
		}

		// append to list
		tags = append(tags, tagModel)
	}

	return tags
}

func getTagCloudEntryLevel(numberOfChilds, minNumberOfChilds, maxNumberOfChilds, levelCount int) int {

	// check the number of childs for negative numbers
	if numberOfChilds < 1 {
		panic(fmt.Sprintf("The number of childs '%v' cannot be less than 1.", numberOfChilds))
	}

	// check max boundary
	if numberOfChilds > maxNumberOfChilds {
		panic(fmt.Sprintf("The number of childs '%v' cannot be greater than the maximum number of childs '%v'.", numberOfChilds, maxNumberOfChilds))
	}

	// check min boundary
	if numberOfChilds < minNumberOfChilds {
		panic(fmt.Sprintf("The number of childs '%v' cannot be smaller than the minimum number of childs '%v'.", numberOfChilds, minNumberOfChilds))
	}

	// check the level count
	if levelCount < 0 {
		panic(fmt.Sprintf("The level count must be greater than 0.", levelCount))
	}

	// calculate the ratio between the "number of childs" to the "maximum number" of childs (0, 1]
	ratioNumberOfChildsToMaxNumberOfChilds := float64(numberOfChilds) / float64(maxNumberOfChilds)
	inverseRation := 1 - ratioNumberOfChildsToMaxNumberOfChilds

	// calculate the level
	level := int(math.Floor(inverseRation*float64(levelCount))) + 1

	return level
}

// sort tags by name
func tagsByName(tag1, tag2 *viewmodel.Tag) bool {
	return tag1.Name < tag2.Name
}

// sort tag cloud entries by name
func tagCloudEntriesByName(tagCloudEntry1, tagCloudEntry2 *viewmodel.TagCloudEntry) bool {
	return tagCloudEntry1.Name < tagCloudEntry2.Name
}
