// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/content"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
)

type FileOrchestrator struct {
	*Orchestrator
}

func (orchestrator *FileOrchestrator) GetFileContentProvider(fileRoute route.Route) *content.ContentProvider {
	file := orchestrator.getFile(fileRoute)
	if file == nil {
		orchestrator.logger.Warn("File %q was not found.", fileRoute)
		return nil
	}

	return file.ContentProvider
}

func (orchestrator *FileOrchestrator) GetFile(fileRoute route.Route) (fileModel viewmodel.File, found bool) {
	file := orchestrator.getFile(fileRoute)
	if file == nil {
		orchestrator.logger.Warn("File %q was not found.", fileRoute)
		return fileModel, false
	}

	convertedModel, err := toViewModel(orchestrator.itemPather(), file)
	if err != nil {
		orchestrator.logger.Warn(err.Error())
		return fileModel, false
	}

	return convertedModel, true
}

func (orchestrator *FileOrchestrator) GetFiles(itemRoute route.Route) (files []viewmodel.File) {
	files = make([]viewmodel.File, 0)

	// get the item
	item := orchestrator.getItem(itemRoute)
	if item == nil {
		orchestrator.logger.Warn("Item %q was not found.", itemRoute)
		return files
	}

	for _, file := range item.Files() {
		fileModel, err := toViewModel(orchestrator.itemPather(), file)
		if err != nil {
			orchestrator.logger.Warn(err.Error())
			continue
		}

		files = append(files, fileModel)
	}

	return files
}

func toViewModel(pathProvider paths.Pather, file *model.File) (fileModel viewmodel.File, err error) {

	// mime type
	mimeType, err := file.MimeType()
	if err != nil {
		return fileModel, fmt.Errorf("Unable to determine mime type of file %q. Error: %s", file, err.Error())
	}

	// hash
	hash, err := file.Hash()
	if err != nil {
		return fileModel, fmt.Errorf("Unable to determine hash of file %q. Error: %s", file, err.Error())
	}

	// last modified date
	lastModifiedDate, err := file.LastModified()
	if err != nil {
		return fileModel, fmt.Errorf("Unable to determine the last modified date of file %q. Error: %s", file, err.Error())
	}

	filePath := file.Route().Path()
	fileRoute := file.Route().Value()

	fileModel = viewmodel.File{
		Parent: file.Parent().Value(),
		Path:   pathProvider.Path(filePath),
		Route:  pathProvider.Path(fileRoute),
		Name:   file.Route().LastComponentName(),

		LastModified: lastModifiedDate,
		MimeType:     mimeType,
		Hash:         hash,
	}

	return fileModel, nil
}
