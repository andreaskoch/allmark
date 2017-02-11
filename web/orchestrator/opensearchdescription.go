// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark/web/view/viewmodel"
	"fmt"
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

	addressPrefix := fmt.Sprintf("%s/", hostname)
	pathProvider := orchestrator.absolutePather(addressPrefix)

	descriptionModel := viewmodel.OpenSearchDescription{
		Title:       fmt.Sprintf("%s Search", rootItem.Title),
		Description: rootItem.Description,
		FavIconURL:  pathProvider.Path("theme/favicon.ico"),
		SearchURL:   pathProvider.Path("search?q={searchTerms}"),
		Tags:        strings.Join(rootItem.MetaData.Tags, " "),
	}

	return descriptionModel
}
