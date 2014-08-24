// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
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
	for _, item := range orchestrator.repository.Items() {

		parsedItem := orchestrator.parseItem(item)
		if parsedItem == nil {
			orchestrator.logger.Warn("Cannot parse item %q", item.String())
			continue
		}

		titleModels = append(titleModels, viewmodel.Title{
			Value:  parsedItem.Title,
			Tokens: strings.Split(parsedItem.Title, " "),
			Route:  orchestrator.itemPather().Path(parsedItem.Route().Value()),
		})

	}

	return titleModels
}
