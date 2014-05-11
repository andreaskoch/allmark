// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"strings"
)

func NewOpenSearchDescriptionOrchestrator(itemIndex *index.Index) OpenSearchDescriptionOrchestrator {
	return OpenSearchDescriptionOrchestrator{
		itemIndex: itemIndex,
	}
}

type OpenSearchDescriptionOrchestrator struct {
	itemIndex *index.Index
}

func (orchestrator *OpenSearchDescriptionOrchestrator) GetDescriptionModel(pathProvider paths.Pather) viewmodel.OpenSearchDescription {

	rootItem := orchestrator.itemIndex.Root()
	if rootItem == nil {
		panic("No root item found")
	}

	tags := make([]string, 0)
	if rootItem.MetaData != nil {
		for _, tag := range rootItem.MetaData.Tags {
			tags = append(tags, tag.Name())
		}
	}

	descriptionModel := viewmodel.OpenSearchDescription{
		Title:       fmt.Sprintf("%s Search", rootItem.Title),
		Description: rootItem.Description,
		FavIconUrl:  pathProvider.Path("theme/favicon.ico"),
		SearchUrl:   pathProvider.Path("search?q={searchTerms}"),
		Tags:        strings.Join(tags, " "),
	}

	return descriptionModel
}
