// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

func createMessageMapperFunc(pathProvider *path.Provider) func(item *repository.Item) view.Model {
	return func(item *repository.Item) view.Model {
		return view.Model{
			Path:        pathProvider.GetWebRoute(item),
			Title:       getTitle(item),
			Content:     item.Content,
			LanguageTag: getTwoLetterLanguageCode(item.MetaData.Language),
		}
	}
}

func getTitle(item *repository.Item) string {
	return "Message posted at " + item.MetaData.Date.String()
}
