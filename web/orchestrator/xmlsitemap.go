// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/elWyatt/allmark/common/paths"
	"github.com/elWyatt/allmark/model"
	"github.com/elWyatt/allmark/web/view/viewmodel"
	"fmt"
	"time"
)

type XmlSitemapOrchestrator struct {
	*Orchestrator
}

func (orchestrator *XmlSitemapOrchestrator) GetSitemapEntires(hostname string) []viewmodel.XmlSitemapEntry {

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found")
	}

	zeroTime := time.Time{}

	children := make([]viewmodel.XmlSitemapEntry, 0)
	for _, item := range orchestrator.getAllItems() {

		// skip virtual items
		if item.IsVirtual() {
			continue
		}

		// item location
		addressPrefix := fmt.Sprintf("%s/", hostname)
		pathProvider := orchestrator.absolutePather(addressPrefix)
		location := pathProvider.Path(item.Route().Value())

		// last modified date
		lastModifiedDate := ""
		if item.MetaData.LastModifiedDate != zeroTime {
			lastModifiedDate = item.MetaData.LastModifiedDate.Format("2006-01-02")
		}

		// images
		images := getImageModels(pathProvider, item)

		children = append(children, viewmodel.XmlSitemapEntry{
			Loc:          location,
			LastModified: lastModifiedDate,
			Images:       images,
		})
	}

	return children
}

func getImageModels(pathProvider paths.Pather, item *model.Item) []viewmodel.XmlSitemapEntryImage {
	var imageModels []viewmodel.XmlSitemapEntryImage

	for _, file := range item.Files() {

		// skip all non-image files
		if !model.IsImageFile(file) {
			continue
		}

		// determine the file location
		fileLocation := pathProvider.Path(file.Route().Value())

		// append the image model
		imageModels = append(imageModels, viewmodel.XmlSitemapEntryImage{
			Loc: fileLocation,
		})

	}

	return imageModels
}
