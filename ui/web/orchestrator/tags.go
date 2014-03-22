// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func NewTagsOrchestrator(itemIndex *index.ItemIndex, pathProvider paths.Pather) TagsOrchestrator {
	return TagsOrchestrator{
		itemIndex:    itemIndex,
		pathProvider: pathProvider,
	}
}

type TagsOrchestrator struct {
	itemIndex    *index.ItemIndex
	pathProvider paths.Pather
}

func (orchestrator *TagsOrchestrator) GetTags() []*viewmodel.Tag {

	// items by tag
	itemsByTag := make(map[string][]*viewmodel.Model)
	for _, item := range orchestrator.itemIndex.Items() {

		itemViewModel := &viewmodel.Model{
			Base: getBaseModel(item),
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

		// sort items
		viewmodel.SortModelBy(sortModelsByDateAndRoute).Sort(items)

		// create view model
		tagModel := &viewmodel.Tag{
			Name:   tag,
			Route:  orchestrator.pathProvider.Path(tag),
			Childs: items,
		}

		// append to list
		tags = append(tags, tagModel)
	}

	// sort the tags
	viewmodel.SortTagBy(tagsByName).Sort(tags)

	return tags
}

func (orchestrator *TagsOrchestrator) getItemTags(item *model.Item) []*viewmodel.Tag {

	tags := make([]*viewmodel.Tag, 0)

	// abort if the item has no tags
	if item == nil || item.MetaData == nil {
		return tags
	}

	for _, tag := range item.MetaData.Tags {

		// create view model
		tagModel := &viewmodel.Tag{
			Name:  tag.Name(),
			Route: orchestrator.pathProvider.Path(tag.Name()),
		}

		// append to list
		tags = append(tags, tagModel)
	}

	return tags
}

// sort tags by name
func tagsByName(tag1, tag2 *viewmodel.Tag) bool {
	return tag1.Name < tag2.Name
}
