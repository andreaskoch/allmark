// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"fmt"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
	"math"
)

var (
	// sort tag cloud entries by name
	tagCloudEntriesByName = func(tagCloudEntry1, tagCloudEntry2 *view.TagCloudEntry) bool {
		return tagCloudEntry1.Name < tagCloudEntry2.Name
	}

	// the maximum number of cloud entry levels
	tagCloudEntryLevels = 6
)

func MapTagCloud(tagmap repository.TagMap, tagPath func(tag *repository.Tag) string) view.TagCloud {
	cloud := make(view.TagCloud, 0)

	minNumberOfItems := 1
	maxNumberOfItems := 1

	for tag, items := range tagmap {

		// calculate the number of items per tag
		numberItemsPerTag := len(items)

		// update the maximum number of items per tag
		if numberItemsPerTag > maxNumberOfItems {
			maxNumberOfItems = numberItemsPerTag
		}

		// update the minimum number of items per tag
		if numberItemsPerTag < minNumberOfItems {
			minNumberOfItems = numberItemsPerTag
		}

		// create a new tag cloud entry
		tagCloudEntry := view.TagCloudEntry{
			Name:           tag.Name(),
			Description:    "",
			AbsoluteRoute:  tagPath(&tag),
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
	view.SortTagCloudBy(tagCloudEntriesByName).Sort(cloud)

	return cloud
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
