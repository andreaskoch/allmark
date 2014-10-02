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

	baseModel := viewmodel.Base{
		RepositoryName:        root.Title,
		RepositoryDescription: root.Description,

		Type:    item.Type.String(),
		Route:   pathProvider.Path(item.Route().Value()),
		Level:   item.Route().Level(),
		BaseUrl: GetBaseUrl(item.Route()),
		Alias:   item.MetaData.Alias,

		PrintUrl: GetTypedItemUrl(item.Route(), "print"),
		JsonUrl:  GetTypedItemUrl(item.Route(), "json"),
		RtfUrl:   GetTypedItemUrl(item.Route(), "rtf"),

		Title:       item.Title,
		Description: item.Description,

		LanguageTag:      item.MetaData.Language,
		CreationDate:     item.MetaData.CreationDate.Format("2006-01-02"),
		LastModifiedDate: item.MetaData.LastModifiedDate.Format("2006-01-02"),
	}

	return baseModel

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
