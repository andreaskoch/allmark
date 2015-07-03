// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"allmark.io/modules/common/route"
	"allmark.io/modules/web/view/viewmodel"
	"fmt"
)

type ConversionModelOrchestrator struct {
	*Orchestrator

	fileOrchestrator *FileOrchestrator
}

func (orchestrator *ConversionModelOrchestrator) GetConversionModel(baseUrl string, route route.Route) (model viewmodel.ConversionModel, found bool) {

	// get the root item
	root := orchestrator.rootItem()
	if root == nil {
		return model, false
	}

	// get the requested item
	item := orchestrator.getItem(route)
	if item == nil {
		orchestrator.logger.Info("There was no item for route %q.", route)
		return model, false
	}

	// create the path provider
	rootPathProvider := orchestrator.absolutePather(fmt.Sprintf("%s/", baseUrl))
	itemContentPathProvider := orchestrator.absolutePather(fmt.Sprintf("%s/%s/", baseUrl, item.Route().Value()))

	// convert content
	convertedContent, err := orchestrator.converter.Convert(orchestrator.getItemByAlias, rootPathProvider, itemContentPathProvider, item)
	if err != nil {
		return model, false
	}

	// create a view model
	model = viewmodel.ConversionModel{
		Base:    getBaseModel(root, item, itemContentPathProvider, orchestrator.config),
		Content: convertedContent,

		// files
		Files: orchestrator.fileOrchestrator.GetFiles(route),
	}

	return model, true
}
