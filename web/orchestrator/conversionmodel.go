// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark/common/route"
	"github.com/andreaskoch/allmark/web/view/viewmodel"
	"fmt"
)

type ConversionModelOrchestrator struct {
	*Orchestrator

	fileOrchestrator *FileOrchestrator
}

func (orchestrator *ConversionModelOrchestrator) GetConversionModel(baseURL string, route route.Route) (model viewmodel.ConversionModel, found bool) {

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
	rootPathProvider := orchestrator.absolutePather(fmt.Sprintf("%s/", baseURL))

	// convert content
	convertedContent, err := orchestrator.converter.Convert(orchestrator.getItemByAlias, rootPathProvider, item)
	if err != nil {
		return model, false
	}

	// create a view model
	model = viewmodel.ConversionModel{
		Base:    getBaseModel(root, item, orchestrator.config),
		Content: convertedContent,

		// files
		Files: orchestrator.fileOrchestrator.GetFiles(route),
	}

	return model, true
}
