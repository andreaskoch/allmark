// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"strings"
)

func getBaseModel(root, item *model.Item, pathProvider paths.Pather) viewmodel.Base {

	return viewmodel.Base{
		RepositoryName:        root.Title,
		RepositoryDescription: root.Description,

		Type:    item.Type.String(),
		Route:   pathProvider.Path(item.Route().Value()),
		Level:   item.Route().Level(),
		BaseUrl: GetBaseUrl(item.Route()),

		PrintUrl: GetTypedItemUrl(item.Route(), "print"),
		JsonUrl:  GetTypedItemUrl(item.Route(), "json"),
		RtfUrl:   GetTypedItemUrl(item.Route(), "rtf"),

		Title:       item.Title,
		Description: item.Description,
	}

}

func GetBaseUrl(route route.Route) string {
	url := route.Value()
	if url != "" {
		return "/" + url + "/"
	}

	return "/"
}

func GetTypedItemUrl(route route.Route, urlType string) string {
	itemPath := GetBaseUrl(route)
	itemPath = strings.TrimSuffix(itemPath, "/")

	if len(itemPath) > 0 {
		return fmt.Sprintf("%s.%s", itemPath, urlType)
	}

	return urlType
}
