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

		var mimeType string

		if content := file.ContentProvider(); content != nil {

			// mime type
			if value, err := content.MimeType(); err == nil {
				mimeType = value
			}

		}

		childs = append(childs, viewmodel.File{
			Parent: file.Parent().Value(),
			Path:   pathProvider.Path(file.Route().Path()),
			Route:  pathProvider.Path(file.Route().Value()),
			Name:   file.Route().LastComponentName(),

			MimeType: mimeType,
		})

	}

	return childs
}
