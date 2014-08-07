// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/ui/web/view/viewmodel"
	"strings"
)

func getBaseModel(root, item *model.Item, pathProvider paths.Pather) viewmodel.Base {

	return viewmodel.Base{
		RepositoryName:        root.Title,
		RepositoryDescription: root.Description,

		Type:    item.Type.String(),
		Route:   pathProvider.Path(item.Route().Value()),
		Level:   item.Route().Level(),
		BaseUrl: getBaseUrl(item.Route()),

		PrintUrl: getTypedItemUrl(item, "print"),
		JsonUrl:  getTypedItemUrl(item, "json"),
		RtfUrl:   getTypedItemUrl(item, "rtf"),

		Title:       item.Title,
		Description: item.Description,
	}

}

func getBaseUrl(route *route.Route) string {
	url := route.Value()
	if url != "" {
		return "/" + url + "/"
	}

	return "/"
}

func getTypedItemUrl(item *model.Item, urlType string) string {
	itemPath := getBaseUrl(item.Route())
	itemPath = strings.TrimSuffix(itemPath, "/")
	return fmt.Sprintf("%s.%s", itemPath, urlType)
}
