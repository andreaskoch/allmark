// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeaheadorchestrator

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel/typeaheadviewmodel"
	"strings"
)

func NewTitlesOrchestrator(itemIndex *index.ItemIndex, pathProvider paths.Pather) TitlesOrchestrator {
	return TitlesOrchestrator{
		itemIndex:    itemIndex,
		pathProvider: pathProvider,
	}
}

type TitlesOrchestrator struct {
	itemIndex    *index.ItemIndex
	pathProvider paths.Pather
}

func (orchestrator *TitlesOrchestrator) GetTitles() []typeaheadviewmodel.Title {

	rootItem := orchestrator.itemIndex.Root()
	if rootItem == nil {
		panic("No root item found")
	}

	titleModels := make([]typeaheadviewmodel.Title, 0)
	for _, child := range orchestrator.itemIndex.Items() {

		titleModels = append(titleModels, typeaheadviewmodel.Title{
			Value:  child.Title,
			Tokens: strings.Split(child.Title, " "),
			Route:  orchestrator.pathProvider.Path(child.Route().Value()),
		})

	}

	return titleModels
}
