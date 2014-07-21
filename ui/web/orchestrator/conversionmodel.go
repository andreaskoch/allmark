// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/conversion"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func NewConversionModelOrchestrator(itemIndex *index.Index, converter conversion.Converter) ConversionModelOrchestrator {
	return ConversionModelOrchestrator{
		itemIndex:        itemIndex,
		converter:        converter,
		fileOrchestrator: NewFileOrchestrator(),
	}
}

type ConversionModelOrchestrator struct {
	itemIndex        *index.Index
	converter        conversion.Converter
	fileOrchestrator FileOrchestrator
}

func (orchestrator *ConversionModelOrchestrator) GetConversionModel(pathProvider paths.Pather, item *model.Item) viewmodel.ConversionModel {

	// get the root item
	root := orchestrator.itemIndex.Root()
	if root == nil {
		return viewmodel.ConversionModel{}
	}

	// convert content
	convertedContent, err := orchestrator.converter.Convert(pathProvider, item)
	if err != nil {
		return viewmodel.ConversionModel{}
	}

	// create a view model
	viewModel := viewmodel.ConversionModel{
		Base:    getBaseModel(root, item, pathProvider),
		Content: convertedContent,

		// files
		Files: orchestrator.fileOrchestrator.GetFiles(pathProvider, item),
	}

	return viewModel
}
