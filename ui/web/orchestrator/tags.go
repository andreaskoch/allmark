// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"math"
)

var (

	// the maximum number of cloud entry levels
	tagCloudEntryLevels = 6
)

func NewTagsOrchestrator(itemIndex *index.ItemIndex, tagPathProvider, itemPathProvider paths.Pather) TagsOrchestrator {
	return TagsOrchestrator{
		itemIndex:        itemIndex,
		tagPathProvider:  tagPathProvider,
		itemPathProvider: itemPathProvider,
	}
}

type TagsOrchestrator struct {
	itemIndex        *index.ItemIndex
	tagPathProvider  paths.Pather
	itemPathProvider paths.Pather
}

func (orchestrator *TagsOrchestrator) GetTags() []*viewmodel.Tag {

	// items by tag
	itemsByTag := make(map[string][]*viewmodel.Model)
	for _, item := range orchestrator.itemIndex.Items() {

		itemViewModel := &viewmodel.Model{
			Base: getBaseModel(item, orchestrator.itemPathProvider),
		}

		tags := []model.Tag{}
		if item.MetaData != nil && len(item.MetaData.Tags) > 0 {
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
			Route:  orchestrator.tagPathProvider.Path(tag),
			Childs: items,
		}

		// append to list
		tags = append(tags, tagModel)
	}

	// sort the tags
	viewmodel.SortTagBy(tagsByName).Sort(tags)

	return tags
}

func (orchestrator *TagsOrchestrator) GetItemTags(item *model.Item) []*viewmodel.Tag {

	tags := make([]*viewmodel.Tag, 0)

	// abort if the item has no tags
	if item == nil || item.MetaData == nil {
		return tags
	}

	for _, tag := range item.MetaData.Tags {

		// create view model
		tagModel := &viewmodel.Tag{
			Name:  tag.Name(),
			Route: orchestrator.tagPathProvider.Path(tag.Name()),
		}

		// append to list
		tags = append(tags, tagModel)
	}

	return tags
}

func (orchestrator *TagsOrchestrator) GetTagCloud() *viewmodel.TagCloud {
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
			Route:          orchestrator.tagPathProvider.Path(tag.Name),
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

	return &cloud
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
