// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

func createDocumentMapperFunc(pathProvider *path.Provider, converterFunc func(item *repository.Item) string) func(item *repository.Item) view.Model {

	return func(item *repository.Item) view.Model {
		return view.Model{
			Path:        pathProvider.GetWebRoute(item),
			Title:       item.Title,
			Description: item.Description,
			Content:     converterFunc(item),
			LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
		}
	}

}
