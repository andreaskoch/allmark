// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
)

func getBaseUrlFromItem(route *route.Route) string {
	url := route.Value()
	if url != "" {
		return "/" + url + "/"
	}

	return "/"
}

func getBaseModel(item *model.Item) viewmodel.Base {
	return viewmodel.Base{
		Type:    item.Type.String(),
		Route:   item.Route().Value(),
		Level:   item.Route().Level(),
		BaseUrl: getBaseUrlFromItem(item.Route()),

		Title:       item.Title,
		Description: item.Description,
	}
}

func sortModelsByDateAndRoute(model1, model2 *viewmodel.Model) bool {

	// ascending by route
	if model1.LastModifiedDate != "" && model2.LastModifiedDate != "" {
		return model1.LastModifiedDate > model2.LastModifiedDate
	}

	return model1.Route > model2.Route
}
