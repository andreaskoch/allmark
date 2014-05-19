// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func NewFileOrchestrator() FileOrchestrator {
	return FileOrchestrator{}
}

type FileOrchestrator struct {
}

func (orchestrator *FileOrchestrator) GetFiles(pathProvider paths.Pather, item *model.Item) []viewmodel.File {

	childs := make([]viewmodel.File, 0)
	for _, file := range item.Files() {

		// file location
		location := pathProvider.Path(file.Route().Value())

		childs = append(childs, viewmodel.File{
			Route: location,
		})
	}

	return childs
}
