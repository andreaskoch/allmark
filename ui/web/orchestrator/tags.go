// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func NewTagsOrchestrator(itemIndex *index.ItemIndex) TagsOrchestrator {
	return TagsOrchestrator{
		itemIndex: itemIndex,
	}
}

type TagsOrchestrator struct {
	itemIndex *index.ItemIndex
}

func (orchestrator *TagsOrchestrator) GetTags(pathProvider paths.Pather) []*viewmodel.Tag {

	tags := make([]*viewmodel.Tag, 0)

	return tags
}
