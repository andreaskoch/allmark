// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"strings"
)

type OpenSearchDescriptionOrchestrator struct {
	*Orchestrator
}

func (orchestrator *OpenSearchDescriptionOrchestrator) GetDescriptionModel(hostname string) viewmodel.OpenSearchDescription {

	rootItem := orchestrator.rootItem()
	if rootItem == nil {
		orchestrator.logger.Fatal("No root item found")
	}

	tags := make([]string, 0)
	if rootItem.MetaData != nil {
		for _, tag := range rootItem.MetaData.Tags {
			tags = append(tags, tag.Name())
		}
	}

	addressPrefix := fmt.Sprintf("http://%s/", hostname)
	pathProvider := orchestrator.absolutePather(addressPrefix)

	descriptionModel := viewmodel.OpenSearchDescription{
		Title:       fmt.Sprintf("%s Search", rootItem.Title),
		Description: rootItem.Description,
		FavIconUrl:  pathProvider.Path("theme/favicon.ico"),
		SearchUrl:   pathProvider.Path("search?q={searchTerms}"),
		Tags:        strings.Join(tags, " "),
	}

	return descriptionModel
}
