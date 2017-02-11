// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark/web/view/viewmodel"
	"strings"
)

type TitlesOrchestrator struct {
	*Orchestrator
}

func (orchestrator *TitlesOrchestrator) GetTitles() []viewmodel.Title {

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found")
	}

	titleModels := make([]viewmodel.Title, 0)
	for _, item := range orchestrator.getAllItems() {

		titleModels = append(titleModels, viewmodel.Title{
			Value:  item.Title,
			Tokens: strings.Split(item.Title, " "),
			Route:  orchestrator.itemPather().Path(item.Route().Value()),
		})

	}

	return titleModels
}
